FROM golang:1.16.5-alpine3.13 AS build

WORKDIR /build

RUN apk add --update make gcc libc-dev git
COPY . .
RUN go build -o /build/tgmailing .


FROM alpine:3.13 AS runtime

RUN apk add --no-cache ca-certificates curl
COPY --from=build /build/tgmailing /tgmailing

EXPOSE 9090

ENTRYPOINT ["/tgmailing"]
