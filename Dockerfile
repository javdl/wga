ARG GO_VERSION=1.22.2
FROM oven/bun:alpine as bun-builder

RUN apk add git
WORKDIR /app/src
COPY . .
RUN bun install
RUN bun run build

FROM golang:${GO_VERSION}-alpine as go-builder

RUN echo "Building with Go version ${GO_VERSION}"

WORKDIR /app/src
COPY --from=bun-builder /app/src /app/src
RUN go mod download && go mod verify
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go build -v -o /run-app .


FROM alpine:latest

COPY --from=go-builder /run-app /usr/local/bin/
CMD ["run-app", "serve", "--http", "0.0.0.0:8090"]
