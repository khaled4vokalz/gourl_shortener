FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

RUN mkdir -p /etc/shortener/

COPY . .

RUN go build -o bin/gourl_shortener cmd/gourl_shortener/main.go
RUN cp bin/gourl_shortener /usr/local/bin/
RUN cp configuration/prod.yaml /etc/shortener/
ENV CONFIG_PATH=/etc/shortener/prod.yaml

# remove all the sources
RUN rm -rf *

EXPOSE 8080

CMD ["/usr/local/bin/gourl_shortener"]
