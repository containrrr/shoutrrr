FROM scratch

COPY ca-certificates.pem /etc/ssl/certs/
COPY shoutrrr /

ENTRYPOINT ["./shoutrrr"]
