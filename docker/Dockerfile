FROM alpine:latest

# Set a working directory
WORKDIR /app

# Copy the entire project
COPY ./dist ./bin
COPY ./config ./config

# Set environment variables
ENV HTTP_PORT=8080
ENV MONGO_SERVICE_URI ""
ENV MONGO_SERVICE_HOST mongo
ENV MONGO_SERVICE_PORT 27017
ENV MONGO_DB app

# Expose the port
EXPOSE "8080:8080"

# Specify the command to run
CMD ["./bin/main"]


