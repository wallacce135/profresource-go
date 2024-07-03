FROM golang:1.20-alpine as builder

WORKDIR /workspace
RUN apk --no-cache add bash gcc fish

COPY . .
RUN go mod download
RUN go build cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /workspace/main .

CMD ["./main"]