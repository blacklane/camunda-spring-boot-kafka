# Prepare
FROM golang:1.13-alpine as baseimg

RUN apk --no-cache upgrade && apk --no-cache add git gcc libc-dev make curl

# Build
FROM baseimg as builder

# First only download the dependencies, so this layer can be cached before we copy the code
COPY ./go.mod ./go.sum /app/
RUN go get -u golang.org/x/lint/golint

# dots are required for global installation
# RUN go get -u github.com/rubenv/sql-migrate/...

WORKDIR /app/
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build

EXPOSE 8000
CMD ["./worker"]
