# 1. AŞAMA: Derleme
FROM node:20-alpine AS builder
WORKDIR /app

# Sadece bağımlılık tanımlarını al
COPY package*.json ./
# Daha hızlı ve tutarlı bir kurulum için 'ci'
RUN npm ci

# Kodun tamamını kopyala ve build al
COPY . .

# Bu satırı ekle: Dosya listesini build sırasında görelim
RUN ls -l src/

RUN npm run build

# 2. AŞAMA: Servis Etme
FROM nginx:stable-alpine

# Özel Nginx konfigürasyonunu kopyala (Routing desteği için)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Builder aşamasından sadece build klasörünü çek
COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]