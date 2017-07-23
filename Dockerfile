FROM golang:1.8 as builder
WORKDIR /go/src/app
COPY app/ .
RUN go-wrapper download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/app/app .
COPY --from=builder /go/src/app/www/ ./www/
EXPOSE 8080
CMD ["./app"]
