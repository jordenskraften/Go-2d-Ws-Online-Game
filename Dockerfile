FROM golang:alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/app cmd/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app /cmd/
COPY frontend/index.html /frontend/

EXPOSE 8080  

CMD ["/cmd/app"]