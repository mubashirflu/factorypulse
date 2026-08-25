# FROM --platform=linux/amd64 golang:1.26-alpine

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -o main .
# RUN chmod +x main

# EXPOSE 8080

# CMD ["./main"]
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o main .

FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./main

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]