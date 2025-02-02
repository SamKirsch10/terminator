FROM golang:1.21 AS buildenv

WORKDIR /app
COPY go.mod go.sum ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM ubuntu:noble

ENV TZ=America/New_York

RUN apt update && \
    apt install -y zfsutils-linux

COPY --from=buildenv /app/zpool_status /app/

ENTRYPOINT [ "/app/zpool_status" ]