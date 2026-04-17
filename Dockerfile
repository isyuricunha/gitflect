FROM golang:1.22-alpine AS builder

WORKDIR /app

# Enable caching for go mod downloads
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gitflect .

FROM alpine:3.19

# Git and curl are essential. Git-lfs for LFS objects.
RUN apk add --no-cache git git-lfs ca-certificates curl

# Setup LFS
RUN git lfs install

COPY --from=builder /app/gitflect /usr/local/bin/gitflect

ENTRYPOINT ["gitflect"]
