FROM golang:1.22.5

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o api .

CMD [ "./product-service-projector/api" ]