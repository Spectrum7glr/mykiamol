FROM golang:1.22 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ./src /app/
RUN go mod download


# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /numbersapi ./cmd/main.go



FROM gcr.io/distroless/base-debian11 AS build-release-stage


WORKDIR /

COPY --from=build-stage /numbersapi /numbersapi

EXPOSE 80

# USER nonroot:nonroot

ENTRYPOINT ["/numbersapi"]