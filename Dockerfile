FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /usr/src/zadarma-exporter

COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

COPY ./main.go  ./ 
COPY ./api ./api

RUN go build -v -o /usr/local/bin/zadarma-exporter

FROM gcr.io/distroless/static-debian12 as final

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /usr/local/bin/zadarma-exporter /usr/local/bin/zadarma-exporter

EXPOSE 9102

ENTRYPOINT ["/usr/local/bin/zadarma-exporter"]
