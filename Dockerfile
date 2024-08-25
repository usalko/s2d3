FROM golang:1.22.5-alpine AS build

RUN apk add --update --no-cache alpine-sdk bash python3

# copy source files and private repo dep
COPY . /go/src/s2d3/
# COPY ./vendor/ /go/src/s2d3/vendor/

# static build the app
WORKDIR /go/src/s2d3
# RUN go mod init
RUN go mod tidy
RUN go install -tags=musl

RUN go build -tags=musl -tags=dynamic cmd/s2d3/main.go

# SHOW CONTENT FROM BUILD FOLDER
RUN ls -la /go/src/s2d3

# create final image
FROM alpine:3.20.2 AS runtime

COPY --from=build /go/src/s2d3/main /usr/bin/s2d3
COPY --from=build /usr/local /usr/local
# Statistics application
COPY ./nue/.dist/prod /statistics/app
ENV STATISTICS_APPLICATION_FOLDER=/statistics/app

RUN apk --no-cache add \
    curl

ENTRYPOINT ["s2d3", "-a", "0.0.0.0"]
