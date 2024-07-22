FROM golang

WORKDIR /container-product-service-projector

COPY . .

RUN go mod download

RUN go build -o api .

CMD [ "/product-service-projector/api" ]