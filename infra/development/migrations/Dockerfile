FROM migrate/migrate:4

# Migrasyon dosyalarını kopyala
COPY . /migrations

# ENTRYPOINT olarak migrate aracını belirle
ENTRYPOINT ["migrate"]

# CMD kısmında veritabanı URL'sini bir değişken ($DB_URL) olarak tanımlıyoruz
ENTRYPOINT ["sh", "-c", "migrate -path=/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE} up"]