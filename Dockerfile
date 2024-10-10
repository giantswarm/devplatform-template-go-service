FROM scratch
COPY --chmod=755 demo-album-catalog /demo-album-catalog
COPY templates /templates
USER 65534:65534
ENTRYPOINT [ "/demo-album-catalog" ]
LABEL org.opencontainers.image.source=https://github.com/giantswarm/demo-album-catalog
