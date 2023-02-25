FROM golang:1.16.3-alpine3.13 as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates git wget

WORKDIR /api
ADD . /api
RUN go mod download
RUN go build -o api .

FROM alpine:3.13.4

RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates tzdata

COPY --from=builder /api/api .
COPY pkey.pem .
COPY active.en.toml .
COPY active.th.toml .
ADD /configs /configs

EXPOSE 8000

ENTRYPOINT ["/api"]
