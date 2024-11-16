FROM golang:1.17-alpine3.13 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV TZ=Asia/Bangkok

RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates git wget
RUN apk add build-base

WORKDIR /api
ADD . /api

CMD ["/app/main"]
RUN go build -o api

FROM alpine:3.13.4

RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates tzdata

COPY --from=builder /api/api .
COPY pkey.pem .
COPY active.en.toml .
COPY active.th.toml .

EXPOSE 8000

ENTRYPOINT ["/api"]

# docker run -d -p 8080:8080 -p 50000:50000 -v /var/run/docker.sock:/var/run/docker.sock -v jenkins_home:/var/jenkins_home jenkins/jenkins
