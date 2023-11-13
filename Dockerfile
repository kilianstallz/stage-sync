
FROM golang:1.21-buster as BUILD

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o /stage-sync ./cmd/stage-sync

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=BUILD /stage-sync /stage-sync

USER nonroot:nonroot

ENTRYPOINT ["/stage-sync"]
