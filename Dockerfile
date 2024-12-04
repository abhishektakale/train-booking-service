
FROM debian:bullseye-slim

WORKDIR /app

COPY ./out/server .

# Expose the port the app runs on
EXPOSE 7001

# Run the Go binary
CMD ["/app/out/server"]
