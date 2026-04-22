FROM golang:1.22-alpine AS builder
ARG HTTP_PROXY
ARG HTTPS_PROXY
ENV HTTP_PROXY=${HTTP_PROXY}
ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV http_proxy=${HTTP_PROXY}
ENV https_proxy=${HTTPS_PROXY}

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

FROM alpine:3.19
ARG HTTP_PROXY
ARG HTTPS_PROXY
ENV HTTP_PROXY=${HTTP_PROXY}
ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV http_proxy=${HTTP_PROXY}
ENV https_proxy=${HTTPS_PROXY}
RUN apk add --no-cache ca-certificates tzdata openssh-client bash python3 py3-pip
RUN pip3 install --no-cache-dir ansible-core==2.14.18 --break-system-packages
RUN ansible-galaxy collection install 'community.general:5.8.0' 'ansible.posix:1.5.4'

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p /var/log/ksk8s /var/lib/ksk8s

EXPOSE 8080
CMD ["./server"]
