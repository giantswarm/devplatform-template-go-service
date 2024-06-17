FROM scratch
COPY --chmod=755 [[project-name]] /[[project-name]]
USER 65534:65534
ENTRYPOINT [ "/[[project-name]]" ]
