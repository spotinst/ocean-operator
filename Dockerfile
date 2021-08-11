# syntax=docker/dockerfile:1.3

#
#
# Images.
ARG BASE_BUILDER=golang:1.16
ARG BASE_RUNTIME=gcr.io/distroless/static:nonroot

#
#
# Builder.
FROM ${BASE_BUILDER} AS builder

# Set the working directory.
WORKDIR /src

# Copy the go.* files and download the modules before copying the rest of the
# source. This allows BuildKit to cache the modules as it will only rerun these
# steps if the go.* files change.
COPY go.* ./
RUN go mod download

# Linker flags.
ARG GO_LDFLAGS

# Copy source.
COPY . ./

# Allow the build container to cache the Go's compiler cache directory.
# Ref: https://docs.docker.com/develop/develop-images/build_enhancements/
# Ref: https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/experimental.md
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="${GO_LDFLAGS}" \
    -o ocean-operator cmd/ocean-operator/main.go

##
##
## Runtime.
FROM ${BASE_RUNTIME} AS runtime

# Copy from builder.
COPY --from=builder /src/ocean-operator /opt/ocean/bin/
COPY --from=builder /src/LICENSE        /opt/ocean/

# Set user.
USER nonroot:nonroot

# Configure the image to run as an executable.
ENTRYPOINT ["/opt/ocean/bin/ocean-operator"]
