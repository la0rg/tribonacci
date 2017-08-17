FROM golang:1.8.1 as builder
WORKDIR /go/src/github.com/la0rg/tribonacci/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/la0rg/tribonacci/app .
EXPOSE 8080
CMD ["./app"]