# STEP 1: build executable binary

# Pull golang image
FROM golang:1.19-alpine as build

# Additional Label
LABEL maintainer="Titanio Yudista<titanioyudista98@gmail.com>"

# Add a work directoryer
WORKDIR /app
# Install make
RUN apk add --no-cache bash make gcc libc-dev

# Cache and install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Copy app files
COPY . .
RUN cp -rf ./.env.example ./.env
# Build app
RUN go build app/main.go

# Expose port
EXPOSE 7000
CMD ["/app/main"]