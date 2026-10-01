FROM node:20-alpine AS webbuilder

RUN npm config set registry https://registry.npmmirror.com
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm install
COPY web/ .
RUN npm run build

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
COPY --from=webbuilder /web/dist ./web/dist

RUN mkdir -p /app/data/uploads

EXPOSE 8080
VOLUME ["/app/data"]

CMD ["./linktree"]
