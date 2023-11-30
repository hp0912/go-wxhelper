FROM golang:alpine as builder
WORKDIR /builder
COPY . .

#ENV GO111MODULE=on
#ENV GOPROXY=https://goproxy.cn,direct

RUN go version
RUN go mod download && go build -o app
RUN ls -lh && chmod -R +x ./*

FROM code.hyxc1.com/open/alpine:3.16.0 as runner
LABEL org.opencontainers.image.authors="lxh@cxh.cn"

EXPOSE 19099
EXPOSE 8080

WORKDIR /app
COPY --from=builder /builder/app ./app
COPY --from=builder /builder/views ./views
CMD ls -lah && ./app