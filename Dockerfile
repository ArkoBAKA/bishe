FROM node:20-alpine AS frontend-builder
WORKDIR /app/vue
RUN corepack enable
COPY vue/package.json vue/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY vue/ ./
RUN pnpm build

FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/server
ARG TARGETOS
ARG TARGETARCH
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH:-amd64}" go build -trimpath -ldflags="-s -w" -o /out/server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates nginx tzdata tini

COPY --from=frontend-builder /app/vue/dist /usr/share/nginx/html
COPY --from=backend-builder /out/server /app/server/server
COPY server/config.yaml /app/server/config.yaml

RUN mkdir -p /app/server/data/uploads/public \
  && rm -f /etc/nginx/http.d/default.conf \
  && cat > /etc/nginx/http.d/app.conf <<'EOF'
server {
  listen 80;
  server_name _;

  client_max_body_size 100m;

  root /usr/share/nginx/html;
  index index.html;

  location /api/ {
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_pass http://127.0.0.1:8080;
  }

  location / {
    try_files $uri $uri/ /index.html;
  }
}
EOF

RUN cat > /entrypoint.sh <<'EOF'
#!/bin/sh
set -e

cd /app/server
export CONFIG_PATH="${CONFIG_PATH:-/app/server/config.yaml}"

./server &
backend_pid="$!"

nginx -g 'daemon off;' &
nginx_pid="$!"

term() {
  kill -TERM "$nginx_pid" 2>/dev/null || true
  kill -TERM "$backend_pid" 2>/dev/null || true
}

trap term INT TERM

while :; do
  if ! kill -0 "$backend_pid" 2>/dev/null; then
    term
    wait "$backend_pid" 2>/dev/null || true
    exit 1
  fi
  if ! kill -0 "$nginx_pid" 2>/dev/null; then
    term
    wait "$nginx_pid" 2>/dev/null || true
    exit 1
  fi
  sleep 1
done
EOF

RUN chmod +x /entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/sbin/tini","--","/entrypoint.sh"]
