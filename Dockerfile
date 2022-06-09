FROM golang:1.18-alpine AS development

ENV PROJECT_PATH=/xm
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

# RUN apk add --no-cache ca-certificates make git bash alpine-sdk nodejs npm
RUN apk add ca-certificates make git bash protobuf protobuf-dev

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH

WORKDIR $PROJECT_PATH

RUN make

FROM alpine:3.16.0 AS production

RUN apk --no-cache add ca-certificates
COPY --from=development /xm/build/xm /usr/bin/xm
COPY --from=development /xm/config.toml /etc/xm/config.toml

USER nobody:nogroup
ENTRYPOINT ["/usr/bin/xm"]