# pull base image which gives all required tools and libraries
FROM golang:1.25.5-alpine

# create folder where app code will be stored
WORKDIR /app

# copy source code from host machine to container
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build executable binary file
RUN go build -o main .

# run the application
CMD ["./main"]