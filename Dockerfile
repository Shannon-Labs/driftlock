# Multi-stage build for a tiny static binary
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /out/driftlock-api ./api-server/cmd/driftlock-api

FROM alpine:3.19 AS runtime
RUN adduser -D -u 65532 driftlock \
    && apk add --no-cache ca-certificates curl
WORKDIR /app
COPY --from=build /out/driftlock-api /app/driftlock-api
USER driftlock
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/driftlock-api"]
