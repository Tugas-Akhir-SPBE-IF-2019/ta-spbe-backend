#!/usr/bin/env bash
set -e

# We require following commands to be installed first, see README.md.
declare -a requiredCommand=("protoc" "protoc-gen-go" "protoc-gen-go-grpc" "protoc-gen-openapi" "protoc-gen-grpc-gateway")
for c in "${requiredCommand[@]}"; do
  if ! [ -x "$(command -v "$c")" ]; then
    echo "Failed: $c is required, but not installed. See README.md" >&2
    exit 1
  fi
done

protoSrcBasePath="${PWD}/proto"
outBaseDir="${PWD}"/tmp/go
# remove previous generated files
rm -Rf "${outBaseDir}"
# recreate out dir
mkdir -p "${outBaseDir}"

protoCmdPrefix="protoc --proto_path=${protoSrcBasePath}/3rdparty/googleapis --proto_path=${protoSrcBasePath}/3rdparty/gnostic --proto_path=${protoSrcBasePath}"

# GENERATE SERVICES API
servicePath="${protoSrcBasePath}/spbe/service"
for svc in "${servicePath}"/*; do
  echo "${svc}"
  $protoCmdPrefix \
  --go_out="${outBaseDir}" \
  --go-grpc_opt=require_unimplemented_servers=true \
  --go-grpc_out="${outBaseDir}" \
  --grpc-gateway_opt generate_unbound_methods=true \
  --grpc-gateway_out="${outBaseDir}" \
  "${svc}"/*.proto
done

servicePath="${protoSrcBasePath}/spbe/mq"
for svc in "${servicePath}"/*; do
  echo "${svc}"
  $protoCmdPrefix \
  --go_out="${outBaseDir}" \
  "${svc}"/*.proto
done

servicePath="${protoSrcBasePath}/spbe/common"
for svc in "${servicePath}"/*; do
  echo "${svc}"
  $protoCmdPrefix \
  --go_out="${outBaseDir}" \
  "${svc}"/*.proto
done

# MOVE TO IMPORTABLE GO
rm -Rfv ./gen/go
mkdir -p ./gen/go
mv "${outBaseDir}"/github.com/bfi-finance/bfi-protobuf/gen/go/bfi ./gen/go
rm -Rfv "${outBaseDir}"