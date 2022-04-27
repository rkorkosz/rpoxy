FROM golang:alpine as build

WORKDIR /go/src/rpoxy
COPY . .
RUN CGO_ENABLED=0 go build -o rpoxy .

FROM gcr.io/distroless/base

COPY --from=build /go/src/rpoxy/rpoxy /rpoxy

VOLUME ["/letsencrypt"]

EXPOSE 443/tcp

ENTRYPOINT ["/rpoxy"]
