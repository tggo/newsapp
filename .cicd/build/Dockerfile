######## First stage: build the executable #######
FROM golang:1.17.2-alpine AS build-env
LABEL stage=builder

WORKDIR /build/

# All these steps will be cached
#<- COPY go.mod and go.sum files to the workspace
COPY ./go.mod .
COPY ./go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

ADD . /build/

ARG APP_NAME="bosters-news-server"
ARG APP_VERSION
ARG GIT_HASH
ARG BUILD_DATE
ARG BUILD_NUMBER


RUN go build  -o /build/news-server -ldflags " -X main.appName=$APP_NAME -X main.release=$APP_VERSION  -X main.gitHash=$GIT_HASH -X main.buildDate=$BUILD_DATE -X main.buildNumber=$BUILD_NUMBER " ./cmd/server

######## Second stage: from scratch #######
FROM alpine:3.12 AS final

# timezone
ENV TZ=Etc/GMT

RUN apk add --no-cache  tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# secure against running as root
RUN adduser -D -u 10000 jasper
USER jasper

WORKDIR /

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy the Pre-built binary file from the previous stage
COPY --from=build-env /build/news-server /bin

# Run the compiled binary.
ENTRYPOINT ["news-server"]
