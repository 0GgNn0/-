FROM golang:1.26-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o linktree .

FROM alpine:3.21

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata curl

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/linktree .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

RUN mkdir -p /app/data/uploads

EXPOSE 8080
VOLUME ["/app/data"]

CMD ["./linktree"]
