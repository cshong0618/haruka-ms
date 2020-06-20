# Base build image
FROM golang:1.13-alpine AS build_base

# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev

WORKDIR /go/src/user

# Force the compiler to use modules

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.

# Because of how the layer caching system works in Docker, the go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

# This image builds the server
FROM build_base AS server_builder

# Here we copy the rest of the source code
COPY . .

# And compile the project
RUN go install ./cmd/rest

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine AS server
# We add the certificates to be able to verify remote weaviate instances
RUN apk add ca-certificates
# Finally we copy the statically compiled Go binary.
COPY --from=server_builder /go/bin/rest user-api

RUN chmod 777 ./user-api

CMD [ "/bin/sh", "-c", "./user-api" ]