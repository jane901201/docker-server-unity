# syntax=docker/dockerfile:1

FROM golang

# Set destination for COPY
WORKDIR /opt/go-socket-server

# 下載 Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
COPY *.env ./

EXPOSE 8888

# Build
RUN go build -o app main.go server.go user.go mysql.go


# Run
CMD [ "/opt/go-socket-server/app" ]
