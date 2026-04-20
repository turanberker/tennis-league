# 1. AŞAMA: Derleme
FROM node:20-alpine AS builder
WORKDIR /app
RUN npm install env-cmd -g # Geçici olarak global yükleyip deneyebilirsin

ARG BUILD_MODE=build

# Sadece bağımlılık tanımlarını al
COPY package*.json ./
# Daha hızlı ve tutarlı bir kurulum için 'ci'
RUN npm ci

# Kodun tamamını kopyala ve build al
COPY . .

RUN npm run ${BUILD_MODE}

# 2. AŞAMA: Servis Etme
FROM nginx:stable-alpine

# Özel Nginx konfigürasyonunu kopyala (Routing desteği için)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Builder aşamasından sadece build klasörünü çek
COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]