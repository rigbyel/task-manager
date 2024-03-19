FROM debian:stable-slim

WORKDIR /app

COPY config/local.yaml ./config/local.yaml
COPY storage/storage.db ./storage/storage.db

COPY task-manager ./task-manager

CMD ["./task-manager"]
