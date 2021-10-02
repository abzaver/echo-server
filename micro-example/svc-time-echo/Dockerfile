FROM golang:1.12-alpine AS build

RUN apk add --update --no-cache build-base

WORKDIR /src/
COPY main*.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/svc-time-echo

FROM scratch
COPY --from=build /bin/svc-time-echo /bin/svc-time-echo
ENTRYPOINT ["/bin/svc-time-echo"]
