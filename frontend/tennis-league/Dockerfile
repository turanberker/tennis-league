# 1. AŞAMA: Derleme (Build)
FROM node:20-alpine AS builder

WORKDIR /app

# Bağımlılıkları kopyala ve kur (Cache dostu)
COPY package.json package-lock.json ./
RUN npm install

# Kaynak kodları kopyala
COPY . .

# Üretim (Production) build'ini al
RUN npm run build

# 2. AŞAMA: Servis Etme (Run)
FROM nginx:stable-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=builder /app/build /usr/share/nginx/html

# Eğer React Router (Client-side routing) kullanıyorsan Nginx konfigürasyonu gerekebilir
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]