FROM golang:1.20-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY Inter/ Inter/

RUN CGO_ENABLED=0 go build -o /invoice .

FROM alpine:3.21

COPY --from=build /invoice /usr/local/bin/invoice

ENTRYPOINT ["invoice"]
