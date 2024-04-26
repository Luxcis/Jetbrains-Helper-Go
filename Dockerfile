FROM golang:latest AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/app/main .
COPY --from=builder /go/src/app/config.toml .
COPY --from=builder /go/src/app/external ./external
COPY --from=builder /go/src/app/static ./static
COPY --from=builder /go/src/app/templates ./templates
CMD ["./main"]
EXPOSE 8080