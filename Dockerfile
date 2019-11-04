FROM golang as builder
WORKDIR /go/src/tericai/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

FROM alpine:latest
WORKDIR /build/
COPY --from=builder /go/src/tericai/app/app /build
ENTRYPOINT ./app
