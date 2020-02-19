package v1

//go:generate sh -c "${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix ${ROOT_DIR}/deploy/crds -nometadata -o ./crd_bindata.go ${ROOT_DIR}/deploy/crds/ocean.spot.io_*_crd.yaml"
