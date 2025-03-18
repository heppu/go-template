FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY demo /bin/demo
ENTRYPOINT [ "/bin/demo" ]
