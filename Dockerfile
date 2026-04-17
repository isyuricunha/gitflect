FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy all source files
COPY . .

# Since we don't have a local go.sum, we run go mod tidy inside the container
RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o gitflect .

FROM alpine:3.19

# Git and curl are essential. Git-lfs for LFS objects.
RUN apk add --no-cache git git-lfs ca-certificates curl

# Setup LFS
RUN git lfs install

COPY --from=builder /app/gitflect /usr/local/bin/gitflect

ENTRYPOINT ["gitflect"]
