FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM nginx:alpine AS prod
WORKDIR /usr/local/bin

RUN apk --no-cache add ca-certificates
RUN mkdir -p /app/keys
COPY keys /app/keys
COPY nginx.conf /etc/nginx/nginx.conf
COPY scripts/nginx_start.sh /usr/local/bin/start.sh
COPY --from=builder /app/main /usr/local/bin/

RUN chown -R nginx:nginx /usr/local/bin/main && \
    chmod +x /usr/local/bin/main && \
    chown -R nginx:nginx /app && \
    chmod +x /usr/local/bin/start.sh

ENV SERVER_PORT="8080"
ENV DB_PATH=/app/app.db
ENV JWK_PATH=/app/keys

EXPOSE 80

CMD ["/usr/local/bin/start.sh"]