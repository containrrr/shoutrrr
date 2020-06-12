FROM golang:alpine as alpine

WORKDIR /shoutrrr 
COPY . .

RUN CGO_ENABLED=0 go build -o shoutrrr ./cmd/shoutrrr 

FROM scratch
COPY --from=alpine \
  /shoutrrr/shoutrrr /shoutrrr

ENTRYPOINT ["./shoutrrr"]