# Builder Stage
FROM golang:1.26.4-alpine as builder

# Copy scripts and application code
COPY ./scripts /scripts
COPY . /app
COPY ./docs /app/docs

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache update && \
    apk --no-cache add gcc g++ make git

# Install air and other dependencies
RUN go install github.com/cosmtrek/air@v1.49.0

# Download Go modules
RUN go mod download
RUN go mod tidy

# Final Stage
FROM alpine:latest

# Install runtime dependencies including Go
RUN apk --no-cache add ca-certificates bash git

# Set working directory
WORKDIR /app

# Copy files from builder stage
COPY --from=builder /app /app
COPY --from=builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"
COPY --from=builder /go/bin/air /usr/local/bin/air

# Ensure scripts and binary are executable
RUN chmod -R +x /app && \
    chmod +x /usr/local/bin/air

# Verify script files
RUN ls -l /app/scripts

# Expose the necessary port
EXPOSE 3022

# Use root user (adjust if not needed)
USER root

# Set the entrypoint to the run.sh script
ENTRYPOINT ["/app/scripts/run-staging.sh"]