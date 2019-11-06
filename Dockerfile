FROM golang as builder

WORKDIR /go/src/tericai/app

RUN go get github.com/globalsign/mgo
RUN go get github.com/globalsign/mgo/bson
RUN go get github.com/joho/godotenv
RUN go get github.com/satori/go.uuid
RUN go get github.com/streadway/amqp
RUN go get github.com/TonPC64/gomon

COPY . /go/src/tericai/common
FROM tericai/go_common:latest as builder
WORKDIR /go/src/tericai/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

FROM alpine:latest
WORKDIR /build/
COPY --from=builder /go/src/tericai/app/app /build
COPY --from=builder /go/src/tericai/app/.env /build/.env
ENTRYPOINT ./app
