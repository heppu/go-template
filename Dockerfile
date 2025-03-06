FROM scratch
ARG TARGET
COPY ${TARGET} /bin/app
ENTRYPOINT [ "/bin/app" ]
