FROM golang:alpine as build

RUN apk --update add ca-certificates
WORKDIR /go/src/rpoxy
COPY . .
RUN CGO_ENABLED=0 go build -o rpoxy .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/rpoxy/rpoxy /rpoxy

VOLUME ["/letsencrypt"]

EXPOSE 443/tcp

ENTRYPOINT ["/rpoxy"]
