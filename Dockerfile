#Web client
FROM node:20-alpine AS node_build
WORKDIR /src
COPY . .
WORKDIR /src/crates/centaureissi_web/src/web
RUN npm ci && npm run build

# Go application
FROM rust:alpine AS builder
WORKDIR /src
RUN apk add --no-cache musl-dev sqlite-dev
ENV RUSTFLAGS=-Ctarget-feature=-crt-static
COPY . .
COPY --from=node_build /src/crates/centaureissi_web/src/web/build/ /src/crates/centaureissi_web/src/web/build/
RUN cargo install --path crates/centaureissi_server

FROM alpine AS runtime
WORKDIR /app
RUN apk add --no-cache sqlite-libs libgcc
COPY --from=builder /usr/local/cargo/bin/centaureissi_server /app/centaureissi_server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/centaureissi_server"]
