ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine AS build

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /project
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /app .

FROM alpine AS runner
COPY --from=build /app /app

EXPOSE 8000
ENTRYPOINT ["/app"]