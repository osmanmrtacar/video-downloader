# syntax=docker/dockerfile:1
FROM golang:1.25-alpine

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
RUN apk add build-base


RUN go build -o video-downloader .

# yt-dlp install
RUN apk add --no-cache yt-dlp

RUN mkdir videos

EXPOSE 8080
CMD ["./video-downloader"]
