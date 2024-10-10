FROM scratch
COPY --chmod=755 {{.ProjectName}} /{{.ProjectName}}
COPY templates /templates
USER 65534:65534
ENTRYPOINT [ "/{{.ProjectName}}" ]
LABEL org.opencontainers.image.source=https://github.com/{{.RepoOwner}}/{{.ProjectName}}
