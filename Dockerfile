FROM golang:latest

# set current workdir
WORKDIR /app

ARG LOG_DIR=/app/logs
RUN mkdir -p ${LOG_DIR}
ENV LOG_FILE_DIR=${LOG_DIR}/app.log

# build steps
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build main.go
EXPOSE 8080

# execution
VOLUME [${LOG_DIR}]
CMD ["./main"]
