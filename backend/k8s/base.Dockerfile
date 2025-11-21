# Go lang with all dependancies as base
FROM golang:1.25-alpine AS base

WORKDIR /app

# Copy and download modules common to all services
COPY ../go.mod ../go.sum ./
RUN go mod download
