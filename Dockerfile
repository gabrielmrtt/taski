FROM golang:1.25

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN apt update && apt install -y make gcc sqlite3 libsqlite3-dev

RUN go mod download
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

# Ensure proper permissions for the binary directory
RUN mkdir -p /usr/src/app/bin && chmod 755 /usr/src/app/bin

EXPOSE 8080

CMD ["make", "dev"]