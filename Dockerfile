FROM golang:1.20 as builder

LABEL org.opencontainers.image.source=https://github.com/cherryReptile/weather-bot

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o ./build cmd/main.go

RUN rm -rf bootstrap/ cmd/ docker/ domain/ domain/ handlers/ pkg/ storage/

CMD ["/app/build"]