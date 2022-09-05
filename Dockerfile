# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /go/src/

COPY ./ ./service

WORKDIR /go/src/service
RUN ls -la

RUN go mod download

RUN go build -o app ./
RUN ls -la

FROM gcr.io/distroless/base-debian10
WORKDIR /
EXPOSE 8081
COPY --from=build /go/src/service /app
CMD ["./app/app"]