# syntax=docker/dockerfile:1.6
#
# Imagen Docker para el servicio de notificaciones (Twilio / WhatsApp)
# escrito en Go con Gin. Expone el puerto 8081 con tres endpoints:
#   POST /notification-medication
#   POST /notification-approved
#   POST /general-message
#
# Build:
#   docker build -t send-notification-twilio:latest .
#
# Run:
#   docker run --rm -p 8081:8081 --env-file .env send-notification-twilio:latest
#

# =============================================================================
# Stage 1 - Build estático del binario Go
# =============================================================================
FROM golang:1.22-alpine AS builder

# git: requerido por algunos `go get`; ca-certificates: por si compila contra HTTPS
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Cachear dependencias antes de copiar el resto del código
COPY go.mod go.sum ./
RUN go mod download

# Código de la app
COPY . .

# Binario estático (sin libc) -> se puede correr en imágenes mínimas
ENV CGO_ENABLED=0 \
    GOOS=linux

# -s -w: stripping de símbolos para reducir tamaño
RUN go build -trimpath -ldflags="-s -w" -o /out/notifications ./


# =============================================================================
# Stage 2 - Imagen final mínima (alpine)
# =============================================================================
FROM alpine:3.19 AS app

LABEL maintainer="Heap LR" \
      service="sendNotificationTwilio" \
      description="Microservicio Go para envío de notificaciones por WhatsApp via Twilio"

# CA certs (Twilio API es HTTPS), tzdata, wget para healthcheck
RUN apk add --no-cache ca-certificates tzdata wget \
 && addgroup -S app && adduser -S -G app -u 1000 app

ENV TZ=America/Argentina/Buenos_Aires \
    GIN_MODE=release

WORKDIR /app

# Binario compilado
COPY --from=builder /out/notifications /app/notifications

# El código hace godotenv.Load() y aborta si no existe .env.
# Creamos uno vacío para que el arranque no falle: las variables reales
# las inyecta docker-compose vía `env_file` / `environment`.
RUN touch /app/.env \
 && chown -R app:app /app

USER app

EXPOSE 8081

# Healthcheck a nivel TCP/HTTP. Los endpoints son POST, GET / devuelve 404
# pero eso ya prueba que el server está vivo. Aceptamos cualquier respuesta HTTP.
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --timeout=3 --server-response \
        --output-document=/dev/null http://127.0.0.1:8081/ 2>&1 \
        | grep -q "HTTP/" || exit 1

ENTRYPOINT ["/app/notifications"]
