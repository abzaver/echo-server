FROM golang:1.17-alpine AS build

RUN apk add --update --no-cache gcc git build-base

WORKDIR /src/
COPY main*.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/tcp-echo-srv

FROM scratch
COPY --from=build /bin/tcp-echo-srv /bin/tcp-echo-srv
ENTRYPOINT ["/bin/tcp-echo-srv"]
