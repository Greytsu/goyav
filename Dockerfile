FROM golang:1.21.4-alpine

#Set working directory
WORKDIR /app

#Download modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

#Build
#CGO_ENABLED=0 disables usage of C libraries
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-goyav

EXPOSE 8080

# Run
CMD ["/api-goyav"]