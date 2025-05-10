# Build the manager binary
FROM docker.io/golang:1.23 AS builder

#ARG TARGETOS
#ARG TARGETARCH

ARG TARGETOS=linux  # 明确变量默认值，否则可以通过命令构建：ocker build --build-arg TARGETOS=linux --build-arg TARGETARCH=amd64 -t my-image .
ARG TARGETARCH=amd64

ARG GOPROXY=https://goproxy.io,direct
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o kubeants-apiserver main.go

FROM alpine:3.21
# 从构建参数获取日期：docker build --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") -t geray/kubeants-apiserver:v1.3.1 .
ARG BUILD_DATE                          
LABEL maintainer="Geray <geray.zhu@gmail.com>" \
    image.authors="geray" \
    image.description="Application packaged by Geray" \
    image.vendor="VMware, Inc." \
    build.date=${BUILD_DATE}   
WORKDIR /
COPY --from=builder /workspace/kubeants-apiserver .
EXPOSE 8080
EXPOSE 8088
USER 65532:65532

ENTRYPOINT ["/kubeants-apiserver"]