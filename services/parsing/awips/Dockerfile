FROM golang:1.24.1 AS build

# Install GEOS
RUN apt-get update && apt-get install -y libgeos-dev

# Set destination for COPY
WORKDIR /

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /awips

FROM alpine:latest AS release

WORKDIR /

COPY --from=build /awips /awips

ENTRYPOINT [ "/awips", "server", "--minlog", "1" ]