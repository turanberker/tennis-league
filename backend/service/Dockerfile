# 1. AŞAMA: Derleme (Build)
FROM golang:1.25-alpine AS builder



# Git ve sertifikalar gibi gerekli araçları kur
RUN apk add --no-cache git

# Çalışma dizinini ayarla
WORKDIR /app

# Önce bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./

# DOSYA KONTROLÜ İÇİN EKLE:
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

# Uygulamayı derle (Buradaki ARG ile hangi uygulamayı derleyeceğimizi seçebiliriz)
ARG APP_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/${APP_NAME}/main.go

# 2. AŞAMA: Çalıştırma (Run)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Sadece binary dosyasını builder aşamasından kopyala
COPY --from=builder /main .

# Çalıştırma komutu
ENTRYPOINT ["./main"]