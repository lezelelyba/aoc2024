FROM golang:1.25.3 AS builder

WORKDIR /release.web

COPY . .

WORKDIR /release.web/cmd/web

RUN go mod download

ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=${VERSION}'" -o /advent2024.webserver

FROM scratch

COPY --from=builder /advent2024.webserver /advent2024.webserver

COPY --from=builder /release.web/cmd/web/templates /templates
COPY --from=builder /release.web/cmd/web/static /static

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["/advent2024.webserver"]