FROM golang:1.25.3 AS builder

WORKDIR /release.cli

COPY . .

WORKDIR /release.cli/cmd/cli

RUN go mod download

ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=${VERSION}'" -o /advent2024.cli

FROM scratch

COPY --from=builder /advent2024.cli /advent2024.cli

ENTRYPOINT ["/advent2024.cli"]