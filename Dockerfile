FROM golang:1.16 AS build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn/,https://mirrors.aliyun.com/goproxy/,direct
WORKDIR  /release

ADD . .
RUN go mod tidy && go mod vendor
# 编译项目,生成二进制文件
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o go_focus_micro_service main.go

FROM debian:bullseye-slim

ENV LANG C.UTF-8

WORKDIR /data

COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 多阶段构建,将上一阶段构建的结果 二进制文件 copy到此阶段镜像中
COPY --from=build /release/go_focus_micro_service .

# 设置时区
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone \
    && cp /etc/apt/sources.list /etc/apt/sources.list.bak \
    && sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y vim

EXPOSE 7066
# 运行二进制文件
CMD ["./go_focus_micro_service"]
