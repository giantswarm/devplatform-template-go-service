package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	gormprom "gorm.io/plugin/prometheus"
)

type storeType int

const (
	storeMemory storeType = iota
	storeDb
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

type album struct {
	ID     uint32  `json:"id" gorm:"primaryKey"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

type albumStore interface {
	Add(album album) (uint32, error)
	Get(id uint32) (album, error)
	List() ([]album, error)
	Update(album album) error
	Remove(id uint32) error
}

type dbStore struct {
	db *gorm.DB
}

func (s *dbStore) Add(album album) (uint32, error) {
	result := s.db.Create(&album)
	if result.Error != nil {
		return 0, result.Error
	}
	return album.ID, nil
}

func (s *dbStore) Get(id uint32) (album, error) {
	var album album = album{
		ID: id,
	}
	result := s.db.First(&album)
	if result.Error != nil {
		return album, result.Error
	}
	return album, nil
}

func (s *dbStore) List() ([]album, error) {
	var albums []album
	result := s.db.Find(&albums)
	if result.Error != nil {
		return nil, result.Error
	}
	return albums, nil
}

func (s *dbStore) Update(album album) error {
	result := s.db.Save(&album)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *dbStore) Remove(id uint32) error {
	result := s.db.Delete(&album{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func NewDbStore() (albumStore, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")

	if username == "" || password == "" || dbName == "" || dbAddress == "" || dbPort == "" {
		return nil, fmt.Errorf("database credentials are not set in environment variables")
	}

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbAddress, dbPort, dbName)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbAddress, username, password, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "albums.",
			SingularTable: false,
		},
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(5)
	db.Exec("CREATE SCHEMA IF NOT EXISTS albums;")
	if err = db.Use(gormprom.New(gormprom.Config{
		DBName:          dbName, // `DBName` as metrics label
		RefreshInterval: 10,     // refresh metrics interval (default 15 seconds)
		StartServer:     false,  // start http server to expose metrics
		MetricsCollector: []gormprom.MetricsCollector{
			&gormprom.Postgres{VariableNames: []string{"Threads_running"}},
		},
	})); err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&album{}); err != nil {
		return nil, err
	}
	return &dbStore{db: db}, nil
}

type inMemoryStore struct {
	albums map[uint32]album
}

func (s *inMemoryStore) Add(album album) (uint32, error) {
	id := uint32(len(s.albums) + 1)
	album.ID = id
	s.albums[id] = album
	return id, nil
}

func (s *inMemoryStore) Get(id uint32) (album, error) {
	album, ok := s.albums[id]
	if !ok {
		return album, fmt.Errorf("album not found")
	}
	return album, nil
}

func (s *inMemoryStore) List() ([]album, error) {
	vals := make([]album, 0, len(s.albums))

	for _, value := range s.albums {
		vals = append(vals, value)
	}
	return vals, nil
}

func (s *inMemoryStore) Update(album album) error {
	id := album.ID
	_, ok := s.albums[id]
	if !ok {
		return fmt.Errorf("album not found")
	}
	s.albums[id] = album
	return nil
}

func (s *inMemoryStore) Remove(id uint32) error {
	_, ok := s.albums[id]
	if !ok {
		return fmt.Errorf("album not found")
	}
	delete(s.albums, id)
	return nil
}

func NewInMemoryStore() albumStore {
	return &inMemoryStore{
		albums: map[uint32]album{
			1: {ID: 1, Title: "All that you can't leave behind", Artist: "U2", Price: 56.99},
			2: {ID: 2, Title: "A night at the opera", Artist: "Queen", Price: 17.99},
		}}
}

type AlbumsHandler struct {
	store albumStore
}

func NewAlbumsHandler(store albumStore) *AlbumsHandler {
	return &AlbumsHandler{store: store}
}

func (h AlbumsHandler) CreateAlbum(c *gin.Context) {
	var album album
	var err error
	if err = c.ShouldBindJSON(&album); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}
	if _, err = h.store.Add(album); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
	}
	c.Redirect(http.StatusFound, "/albums")
}

func (h AlbumsHandler) ListAlbums(c *gin.Context) {
	res, err := h.store.List()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "albums.html", gin.H{"Albums": res})
}

func (h AlbumsHandler) GetAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}
	res, err := h.store.Get(uint32(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "album.html", gin.H{"Album": res})
}

func (h AlbumsHandler) UpdateAlbum(c *gin.Context) {
	var album album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}

	if err := h.store.Update(album); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/albums")
}

func (h AlbumsHandler) DeleteAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}
	err = h.store.Remove(uint32(id))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.Redirect(http.StatusFound, "/albums")
}

func setupPrometheusExporter() {
	// Create a new Prometheus gauge metric
	heartbeatMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "my_webservice_heartbeat",
		Help: "Web Service heartbeat metric",
	})

	// Set the value of the heartbeat metric to 1
	heartbeatMetric.Set(1)

	// Register the metric with the Prometheus default registry
	prometheus.MustRegister(heartbeatMetric)
}

func ginPrometheusMetrics() gin.HandlerFunc {
	httpDurations := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_durations_histogram_seconds",
		Help:       "HTTP latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
		[]string{"method", "route"})

	prometheus.MustRegister(httpDurations)
	return func(c *gin.Context) {
		t := time.Now()

		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		basePath := strings.Split(c.Request.URL.Path, "/")[1]
		httpDurations.WithLabelValues(c.Request.Method, basePath).Observe(latency.Seconds())
	}
}

func setupRedisCache() (persist.CacheStore, error) {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")
	if redisAddr == "" || redisPort == "" {
		return nil, fmt.Errorf("redis address:port is not set in environment variables")
	}
	return persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%s", redisAddr, redisPort),
		Username: redisUsername,
		Password: redisPassword,
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	})), nil
}

func setupRouter(storeType storeType) *gin.Engine {
	router := gin.New()
	logger, _ := zap.NewProduction()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	setupPrometheusExporter()
	router.Use(ginPrometheusMetrics())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.SetFuncMap(template.FuncMap{
		"Version": func() string { return version },
		"Commit":  func() string { return commit },
		"Date":    func() string { return date },
	})
	router.LoadHTMLGlob("./templates/*")

	var store albumStore
	var cacheStore persist.CacheStore
	if storeType == storeMemory {
		store = NewInMemoryStore()
		cacheStore = persist.NewMemoryStore(5 * time.Second)
	} else {
		dbStore, err := NewDbStore()
		if err != nil {
			log.Panicf("Critical error: %s", err)
		}
		store = dbStore
		cacheStore, err = setupRedisCache()
		if err != nil {
			log.Panicf("Critical error: %s", err)
		}
	}
	albumsHandler := NewAlbumsHandler(store)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong\n")
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/addalbum", func(c *gin.Context) {
		c.HTML(http.StatusOK, "addalbum.html", gin.H{})
	})
	router.GET("/albums", albumsHandler.ListAlbums)
	router.GET("/albums/:id", albumsHandler.GetAlbum)
	router.GET("/cached/albums", cache.CacheByRequestURI(cacheStore, 10*time.Second), albumsHandler.ListAlbums)
	router.GET("/cached/albums/:id", cache.CacheByRequestURI(cacheStore, 10*time.Second), albumsHandler.GetAlbum)
	router.POST("/albums", albumsHandler.CreateAlbum)
	router.PUT("/albums/:id", albumsHandler.UpdateAlbum)
	router.DELETE("/albums/:id", albumsHandler.DeleteAlbum)

	return router
}

func main() {
	log.Printf("Web service starting, version: '%s', commit: '%s', build date: '%s', built by: '%s'", version, commit, date, builtBy)
	listeningPort := os.Getenv("LISTEN_PORT")
	if listeningPort == "" {
		listeningPort = "8080"
	}
	memoryStoreMode, found := os.LookupEnv("IN_MEMORY_STORE")

	var router *gin.Engine
	if found && strings.ToLower(memoryStoreMode) == "true" {
		router = setupRouter(storeMemory)
	} else {
		router = setupRouter(storeDb)
	}
	// Listen and Server on the LocalHost:Port
	err := router.Run(":" + listeningPort)
	if err != nil {
		log.Panicf("Critical error: %s", err)
	}
}
