FROM golang:1.21 AS buildenv

WORKDIR /app
COPY go.mod go.sum ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM debian

ENV TZ=America/New_York

COPY debian.sources /etc/apt/sources.list.d/debian.sources

RUN apt update && \
    apt -t experimental install -y zfsutils-linux

COPY --from=buildenv /app/zpool_status /app/

ENTRYPOINT [ "/app/zpool_status" ]
