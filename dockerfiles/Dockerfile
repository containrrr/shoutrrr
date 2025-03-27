FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c AS alpine

RUN apk add --no-cache ca-certificates

FROM scratch

COPY --from=alpine \
    /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt
COPY shoutrrr/shoutrrr /

ENTRYPOINT ["/shoutrrr"]
