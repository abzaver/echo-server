# Dockerfile for User Service
FROM golang:1.12-alpine AS build

RUN apk add git --update --no-cache build-base

WORKDIR /src/
COPY main*.go go.* /src/

# update all modules listed in go.mod
RUN go get -u
# adds any missing imported modules and remove unnecessary ones
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /bin/user-service

FROM scratch
COPY --from=build /bin/user-service /bin/user-service

EXPOSE 8000
ENTRYPOINT ["/bin/user-service"]
