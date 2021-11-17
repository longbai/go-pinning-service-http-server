FROM golang:1.10 AS build
WORKDIR /go/src
COPY service ./go
COPY main.go .

ENV CGO_ENABLED=0
RUN service get -d -v ./...

RUN service build -a -installsuffix cgo -o openapi .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /service/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
