FROM golang:1.25-alpine AS base

# Install tools needed for building Go programs
RUN apk add --no-cache git build-base

# All services will build inside this directory
WORKDIR /app

# Copy module files so Go dependencies can be downloaded once and cached
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
