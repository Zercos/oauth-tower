FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .


FROM nginx:latest

COPY nginx.conf /etc/nginx/nginx.conf
WORKDIR /app
COPY . .
ENV SERVER_PORT="8080"
ENV DB_PATH=/app/app.db
ENV JWK_PATH=/app/keys

WORKDIR /usr/share/nginx/html
COPY --from=builder /app/main /usr/share/nginx/html/main

EXPOSE 80
