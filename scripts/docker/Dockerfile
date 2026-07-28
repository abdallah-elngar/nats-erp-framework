# مرحلة البناء
FROM golang:1.21-alpine AS builder

WORKDIR /app

# نسخ ملفات الاعتماديات
COPY go.mod go.sum ./
RUN go mod download

# نسخ الكود المصدري
COPY . .

# بناء التطبيق
RUN CGO_ENABLED=0 GOOS=linux go build -o nats ./cmd/nats/main.go

# مرحلة التشغيل
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# نسخ الملفات من مرحلة البناء
COPY --from=builder /app/nats .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/apps ./apps
COPY --from=builder /app/config ./config

# المنفذ
EXPOSE 8080

# تشغيل التطبيق
ENTRYPOINT ["./nats"]
CMD ["serve", "--host=0.0.0.0", "--port=8080"]