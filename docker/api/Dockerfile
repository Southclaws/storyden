FROM golang:1.23

RUN git config --global --add safe.directory /server

WORKDIR /server
ADD . .

# Install Taskfile
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Build the API server binary
RUN task backend

ENTRYPOINT [ "task", "production" ]
