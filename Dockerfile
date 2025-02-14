FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /url_shortener ./cmd

EXPOSE 8000

CMD ["/url_shortener"] 