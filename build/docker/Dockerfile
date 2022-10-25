## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg

RUN go build -o /golinks ./cmd/golinks/golinks.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /golinks /golinks

EXPOSE 8080

USER nonroot:nonroot

COPY --from=busybox:1.35.0-uclibc /bin /bin

ENTRYPOINT ["/golinks"]