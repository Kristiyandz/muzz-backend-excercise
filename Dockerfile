# syntax=docker/dockerfile:1

FROM golang:1.22.0

# set destination for COPY
WORKDIR /app

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY *.go ./
COPY . ./

# build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /muzz-docker-api

# expose port 8080
EXPOSE 8080

# run the binary
CMD [ "/muzz-docker-api" ]

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
#COPY *.go ./
COPY . ./