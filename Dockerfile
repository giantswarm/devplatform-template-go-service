FROM scratch
COPY --chmod=755 devplatform-template-go-service /devplatform-template-go-service
USER 65534:65534
ENTRYPOINT [ "/devplatform-template-go-service" ]
