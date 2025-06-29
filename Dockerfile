FROM golang:1.23.8-bullseye AS builder
WORKDIR /movieSpider
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN CGO_ENABLED=0 go build -o movieSpider

FROM  hairyhenderson/upx AS upx
WORKDIR /movieSpider
COPY --from=builder /movieSpider/movieSpider .
RUN upx movieSpider

FROM frolvlad/alpine-glibc
WORKDIR /app
ENV PATH=/app:$PATH
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
    apk add -U tzdata --no-cache &&\
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime &&\
    echo $TZ > /etc/timezone
COPY --from=upx /movieSpider/movieSpider .
ENTRYPOINT ["./movieSpider"]
CMD ["-f","config.yaml"]