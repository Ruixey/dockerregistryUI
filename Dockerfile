FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o registryUI


FROM alpine

WORKDIR /app
COPY --from=builder /app/registryUI .
COPY templates templates
COPY static static

ENV REGISTRY_HOST ""
ENV REGISTRY_PROTOCOL https
ENV REGISTRY_URL ""
ENV IGNORE_INSECURE_HTTPS false
ENV REGISTRY_BASIC_AUTH_USER ""
ENV REGISTRY_BASIC_AUTH_PASSWORD ""

EXPOSE 8080

ENTRYPOINT ["./registryUI"]
