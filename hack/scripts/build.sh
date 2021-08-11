#!/usr/bin/env bash
set -euo pipefail

# Version.
VERSION="${VERSION:-$(cat internal/version/VERSION)}"

# Revision.
REVISION="${REVISION:-$(git rev-parse --short HEAD)}"

# Release channel the image is being built for. Defaults to `stable`.
RELEASE_CHANNEL="${RELEASE_CHANNEL:-stable}"

# Target platform the image is being built for. Defaults to `linux/amd64`.
TARGET_PLATFORM="${TARGET_PLATFORM:-linux/amd64}"

# Image registry. Defaults to `docker.io`.
IMAGE_REGISTRY="${IMAGE_REGISTRY:-docker.io}"

# Image repository. Defaults to `spotinst`.
IMAGE_REPOSITORY="${IMAGE_REPOSITORY:-spotinst}"

# Image name. Defaults to `ocean-operator`.
IMAGE_NAME="${IMAGE_NAME:-ocean-operator}"

# Image tag. One of `$VERSION`, `$VERSION-$RELEASE_CHANNEL`, `latest`.
IMAGE_TAG="${IMAGE_TAG:-${VERSION}}"

# Image output. One of `push` (push image into registry), `load` (load image into Docker).
IMAGE_OUTPUT="${IMAGE_OUTPUT:-push}"

function log() {
	nanosec="$(date +%N)"
	timestamp="$(date -u +"%Y-%m-%dT%H:%M:%S").${nanosec:0:3}Z"
	log_level="$1"
	shift
	message="$*"
	echo -e "${timestamp} [${log_level}] ${message}" >&2
}

function log_info() {
	for message; do log "INFO" "${message}"; done
}

function log_fatal() {
	while (($# > 1)); do
		log "FATAL" "${1}"
		shift
	done
	exit "${1}"
}

function lowercase() {
	echo -n "${1}" | tr '[:upper:]' '[:lower:]'
}

function trim() {
	echo -n "${1}" | xargs
}

function docker_buildx() {
	# Enable extended build capabilities with BuildKit.
	cmd_args="$*"
	cmd="DOCKER_CLI_EXPERIMENTAL=enabled docker buildx ${cmd_args}"
	log_info "executing: ${cmd}"
	eval "${cmd}"
	echo $?
}

function docker_ensure_builder() {
	log_info "ensuring availability of buildkit builder"
	builder_name="ocean-operator-builder"
	builder_status="$(docker_buildx "inspect ${builder_name} >/dev/null 2>&1")"
	[[ "${builder_status}" -gt 0 ]] &&
		docker_buildx "create --name ${builder_name} --use --buildkitd-flags=--debug"
	docker_buildx "inspect ${builder_name} --bootstrap"
}

function docker_ensure_platforms() {
  log_info "ensuring platforms"
  platforms="$(trim $(lowercase "${TARGET_PLATFORM}"))"
	case "${platforms}" in
	"all")
		echo "linux/amd64,linux/arm64"
		;;
	"linux/amd64" | "linux/arm64")
		echo "${platforms}"
		;;
	*)
		log_fatal "unsupported platform: ${platforms}" 128
		;;
	esac
}

function docker_ensure_image() {
	image_registry="$(trim $(lowercase "${IMAGE_REGISTRY}"))"
	image_repository="$(trim $(lowercase "${IMAGE_REPOSITORY}"))"
	image_name="$(trim $(lowercase "${IMAGE_NAME}"))"
	image_tag="$(trim $(lowercase ${IMAGE_TAG}))"
	release_channel="$(trim $(lowercase "${RELEASE_CHANNEL}"))"
	[[ "${image_tag}" != "latest" && "${release_channel}" != "stable" ]] &&
		image_tag+="-${release_channel}"
	echo "${image_registry}/${image_repository}/${image_name}:${image_tag}"
}

function docker_build() {
	image="$(docker_ensure_image)"
	platforms="$(docker_ensure_platforms)"

	args=(
		"build"
		"--platform ${platforms}"
		"--tag ${image}"
		"--label org.opencontainers.image.version=${VERSION}"
		"--label org.opencontainers.image.revision=${REVISION}"
		"--label org.opencontainers.image.created=$(date --rfc-3339=seconds | sed "s/ /T/")"
		"--label org.opencontainers.image.authors=spotinst"
		"--${IMAGE_OUTPUT}"
		"."
	)

	log_info "building image: ${image} (${platforms})"
	docker_buildx "${args[@]}"
}

function main() {
	docker_ensure_platforms
	docker_ensure_builder
	docker_build
}

main "$@"
