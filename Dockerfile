FROM golang:1.15.5-alpine3.12 AS build

WORKDIR /build

RUN apk add --update make gcc libc-dev git

COPY . .

WORKDIR /build/ads/tgmailing

RUN go build -o /build/tgmailing .


FROM alpine:3.11 AS runtime

RUN apk add --no-cache ca-certificates curl
COPY --from=build /build/tgmailing /tgmailing

ENTRYPOINT ["/tgmailing"]
