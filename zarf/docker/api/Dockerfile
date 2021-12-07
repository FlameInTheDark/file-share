# Build the Go Binary.
FROM golang:1.16.1 as build_api
ENV CGO_ENABLED 0
ARG VCS_REF

  # Create a location in the container for the source code. Using the
  # default GOPATH location.
RUN mkdir -p /service

  # Copy the source code into the container.
WORKDIR /service
COPY .. .

WORKDIR /service/app/api
RUN go build -ldflags "-X main.build=${VCS_REF}"

# Run the Go Binary in Alpine.
FROM alpine:3.13
ARG BUILD_DATE
ARG VCS_REF

# Installing curl to use healthchecks
RUN apk add curl

COPY --from=build_api /service/app/api/api /service/api
WORKDIR /service
CMD ["./api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
org.opencontainers.image.title="file-api" \
org.opencontainers.image.authors="" \
org.opencontainers.image.source="https://github.com/FlameInTheDark/file-share" \
org.opencontainers.image.revision="${VCS_REF}" \
org.opencontainers.image.vendor="FlameInTheDark"
