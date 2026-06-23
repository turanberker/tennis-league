# ========================================================
# 1. AŞAMA: DERLEME (BUILDER)
# ========================================================
FROM golang:1.25-alpine AS builder

# Git ve build araçlarını yükle (gerekirse)
RUN apk add --no-cache git

WORKDIR /app

# Önce go.work ve tüm modüllerin go.mod / go.sum dosyalarını kopyala (Cache avantajı için)
COPY go.work go.work.sum* ./
COPY common/consumer/go.mod ./common/consumer/
COPY common/http/go.mod ./common/http/
COPY common/lib/go.mod ./common/lib/
COPY common/security/go.mod ./common/security/
COPY service/main/go.mod ./service/main/
COPY service/user/go.mod ./service/user/
COPY service/user-interface/go.mod ./service/user-interface/

# Bağımlılıkları indir
RUN go mod download

# Tüm kaynak kodları kopyala
COPY . .

# Uygulamaları statik olarak derle CGO_ENABLED=0
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main-api ./service/main/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main-consumer ./service/main/cmd/consumer
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main-publisher ./service/main/cmd/publisher
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-api ./service/user/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-consumer ./service/user/cmd/consumer

# ========================================================
# 2. AŞAMA: ÇALIŞTIRMA (RUNNER) - Her uygulama için ayrı hedef
# ========================================================

# --- service/main/api için ---
FROM alpine:latest AS main-api
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bin/main-api .
CMD ["./main-api"]

# --- service/main/consumer için ---
FROM alpine:latest AS main-consumer
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bin/main-consumer .
CMD ["./main-consumer"]

# --- service/main/publisher için ---
FROM alpine:latest AS main-publisher
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bin/main-publisher .
CMD ["./main-publisher"]

# --- service/user/api için ---
FROM alpine:latest AS user-api
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bin/user-api .
CMD ["./user-api"]

# --- service/user/api için ---
FROM alpine:latest AS user-consumer
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bin/user-consumer .
CMD ["./user-consumer"]

