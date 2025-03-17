FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ARG TARGET
COPY ${TARGET} /bin/app
ENTRYPOINT [ "/bin/app" ]
