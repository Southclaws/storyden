FROM golang:latest

WORKDIR /server
ADD . .

# Node.js
RUN apt update && \
    apt install -y nodejs npm && \
    rm -rf /var/lib/apt/lists/*

# Install prisma command for automatic migrations.
RUN npm install --global prisma

# Install Taskfile
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Install prisma client code generation tool and generate prisma bindings
RUN task generate

# Build the API server binary
RUN task api

ENTRYPOINT [ "task", "production" ]
