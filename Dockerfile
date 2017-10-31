FROM golang:1.8-alpine

RUN apk --no-cache add git g++

RUN go get github.com/sheng/air && \
    go get github.com/astaxie/beego/orm && \
    go get github.com/mattn/go-sqlite3 && \
    go get github.com/kamildrazkiewicz/go-flow

WORKDIR /go/src/app

COPY ["config.toml", "keymaker.go", "/go/src/app/"]
