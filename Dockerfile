FROM golang:latest AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/app/main .
CMD ["./main"]
EXPOSE 8080