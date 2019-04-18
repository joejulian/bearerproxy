FROM golang:1.12 as build

RUN CGO_ENABLED=0 GOOS=linux go get github.com/joejulian/bearerproxy

FROM scratch

COPY --from=build /go/bin/bearerproxy /

CMD ["/bearerproxy"]
