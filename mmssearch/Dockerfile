
FROM golang:latest
WORKDIR /go/src/mmssearch
COPY main.go /go/src/mmssearch
COPY vendor /go/src/mmssearch/vendor
COPY static /go/src/mmssearch/static
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl
WORKDIR /root/
COPY --from=0 go/src/mmssearch/main .
COPY --from=0 go/src/mmssearch/static static
CMD ["./main"]
LABEL version=demo-3
