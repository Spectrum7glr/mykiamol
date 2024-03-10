FROM golang:1.22 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ./src /app/
RUN go mod download


# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./numbersweb ./cmd/main.go


# for debugging
FROM gcr.io/distroless/base-debian11:debug  AS build-release-stage

# FROM gcr.io/distroless/base-debian11  AS build-release-stage

# ENV API_ENDPOINT=http://localhost:8091/rnd
WORKDIR /app

COPY --from=build-stage /app/numbersweb /app/numbersweb
# COPY --from=build-stage /app/cmd/default.env /app/default.env

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/app/numbersweb"]