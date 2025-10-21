# Multi-stage build for a tiny static binary
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /out/driftlock-api ./api-server/cmd/driftlock-api

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/driftlock-api /driftlock-api
USER nonroot
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/driftlock-api"]
