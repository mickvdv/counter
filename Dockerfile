# BUILD
FROM golang:1.24.2-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN go build -o app .

FROM golang:1.24.2-alpine AS run
# Set the working directory inside the container
WORKDIR /app

EXPOSE 8080

# Copy the built application from the build stage
COPY --from=build /app/app .

# Command to run the application
ENTRYPOINT ["./app"]

