FROM golang:latest

ENV GO111MODULE=on
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY ./ ./
RUN go build -o main ./cmd/main.go
EXPOSE 8080
CMD ["./main"]