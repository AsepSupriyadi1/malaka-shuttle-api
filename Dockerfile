# Gunakan base image golang
FROM golang:1.23.4

# Set working directory
WORKDIR /app

# Copy semua file
COPY . .

# Download dependency
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Port container
EXPOSE 8080

# Jalankan binary
CMD ["./main"]
