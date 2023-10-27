FROM golang:1.21.3-alpine3.17 as build
WORKDIR /src
COPY go.* ./
RUN go mod download

COPY internal internal
COPY cmd cmd
RUN go build -o controller ./cmd

FROM alpine:3.17
WORKDIR /opt/controller
COPY --from=build /src/controller .
CMD [ "./controller" ]