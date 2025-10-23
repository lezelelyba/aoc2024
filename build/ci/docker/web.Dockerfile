FROM golang:1.25.3 AS builder

WORKDIR /release.web

COPY . .

WORKDIR /release.web/cmd/web

# update doc files
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency

# update golang documentation
RUN go install github.com/navid-m/arrow@latest
RUN mkdir godocs && \
    cd godocs && \
    arrow /release.web/ --name="AoC 2024 Solver" && \
    mv docs/* . && \
    rmdir docs && \
    cd ..

RUN go mod download

ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=${VERSION}'" -o /advent2024.webserver

FROM scratch

COPY --from=builder /advent2024.webserver /advent2024.webserver

COPY --from=builder /release.web/cmd/web/templates /templates
COPY --from=builder /release.web/cmd/web/static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /release.web/cmd/web/godocs /godocs



EXPOSE 8080

CMD ["/advent2024.webserver"]