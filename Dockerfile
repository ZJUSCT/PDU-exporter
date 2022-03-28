FROM golang:1.16 as builder
WORKDIR /app
ENV GO111MODULE=on GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o PDU-exporter .

FROM scratch
COPY --from=builder /app/PDU-exporter /PDU-exporter
ENTRYPOINT ["/PDU-exporter"]
CMD ["-c", "/config.yml"]
