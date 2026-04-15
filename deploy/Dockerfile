# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/web/admin
COPY web/admin/package*.json ./
RUN npm ci
COPY web/admin/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.25-alpine AS backend
WORKDIR /app
RUN apk add --no-cache gcc musl-dev git ca-certificates
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
RUN go mod download -x
COPY . .
COPY --from=frontend /app/web/admin/dist ./web/admin/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./cmd/server && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/worker ./cmd/worker

# Stage 3: Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/worker .
COPY config/config.yaml ./config/
COPY config/rbac_model.conf ./config/

EXPOSE 8080
CMD ["./server"]
