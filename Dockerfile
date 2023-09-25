FROM golang:alpine as builder
WORKDIR /builder
COPY . .

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go version
RUN go mod download && go build -o app
RUN ls -lh && chmod +x ./app

FROM repo.lxh.io/alpine:3.16.0 as runner
LABEL org.opencontainers.image.authors="lxh@cxh.cn"

EXPOSE 19099

WORKDIR /app
COPY --from=builder /builder/app ./app
CMD ./app