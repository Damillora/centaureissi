#Web client
FROM node:20-alpine AS node_build
WORKDIR /src
COPY . .
WORKDIR /src/pkg/web
RUN npm ci && npm run build

# Go application
FROM golang:1.24-alpine AS build
WORKDIR /go/src/centaureissi
COPY . .
# COPY --from=node_build /src/pkg/web/build/ /go/src/centaureissi/pkg/web/build/
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /centaureissi -ldflags '-extldflags "-static"' -tags timetzdata github.com/Damillora/centaureissi/cmd/centaureissi

FROM scratch AS runtime
WORKDIR /app
COPY --from=build /centaureissi /app/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/centaureissi"]
