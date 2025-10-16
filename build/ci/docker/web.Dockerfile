FROM golang:1.25.3 AS builder

WORKDIR /release.web

COPY . .

WORKDIR /release.web/cmd/web

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /advent2024.webserver


FROM scratch

COPY --from=builder /advent2024.webserver /advent2024.webserver

COPY --from=builder /release.web/cmd/web/templates /templates
COPY --from=builder /release.web/cmd/web/static /static

EXPOSE 8080

CMD ["/advent2024.webserver"]