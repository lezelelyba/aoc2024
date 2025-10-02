FROM golang:1.24.7 AS builder

WORKDIR /release.cli

COPY . .

WORKDIR /release.cli/cmd/cli

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /advent2024.cli

FROM scratch

COPY --from=builder /advent2024.cli /advent2024.cli

CMD ["/advent2024.cli"]