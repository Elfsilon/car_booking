FROM golang:1.22.4-bullseye

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./bin/app ./cmd/main.go

CMD [ "./bin/app" ]