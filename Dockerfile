FROM golang:1.24.0-alpine AS builder
WORKDIR /movieSpider
ENV GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 \
    GOPATH=/root/gopath \
    GOPROXY=https://goproxy.cn,direct \
    GO111MODULE='on' \
    GIT_TERMINAL_PROMPT=1
COPY . .
RUN CGO_ENABLED=0 go build -o movieSpider

FROM  hairyhenderson/upx AS upx
WORKDIR /movieSpider
COPY --from=builder /movieSpider/movieSpider .
RUN upx movieSpider

FROM frolvlad/alpine-glibc
WORKDIR /app
ENV PATH=/app:$PATH TZ=Asia/Shanghai
RUN apk add -U tzdata --no-cache &&\
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime &&\
    echo $TZ > /etc/timezone
COPY --from=upx /movieSpider/movieSpider .
ENTRYPOINT ["./movieSpider"]
CMD ["-f","config.yaml"]