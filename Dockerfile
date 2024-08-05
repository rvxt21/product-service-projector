FROM golang:1.22.5

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd/api

RUN go build -o api .

CMD [ "./api" ]