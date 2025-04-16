FROM golang:1.24.2 AS build

# Install GEOS
RUN apt-get update && apt-get install -y libgeos-dev

# Set destination for COPY
WORKDIR /

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o /awips

# Use glibc-based image in the final stage
FROM debian:bookworm-slim AS release

WORKDIR /

# Install runtime dependencies
RUN apt-get update && apt-get install -y libgeos-c1v5

COPY --from=build /awips /awips

ENTRYPOINT [ "/awips", "server", "--minlog", "1" ]
