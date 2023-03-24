#!/usr/bin/env bash
set -e

# We require following commands to be installed first, see README.md.
declare -a requiredCommand=("protoc" "protoc-gen-openapi" "protoc-gen-doc")
for c in "${requiredCommand[@]}"; do
  if ! [ -x "$(command -v "$c")" ]; then
    echo "Failed: $c is required, but not installed. See README.md" >&2
    exit 1
  fi
done

protoSrcBasePath="${PWD}/proto"
outBaseDir="${PWD}"/tmp/docs
# remove previous generated files
rm -Rf "${outBaseDir}"
# recreate out dir
mkdir -p "${outBaseDir}"

protoCmdPrefix="protoc --proto_path=${protoSrcBasePath}/3rdparty/googleapis --proto_path=${protoSrcBasePath}/3rdparty/gnostic --proto_path=${protoSrcBasePath}"

# GENERATE SERVICES API
servicePath="${protoSrcBasePath}/spbe/service"
for svc in "${servicePath}"/*; do
  echo "${svc}"
  docDir="${outBaseDir}"/spbe/service/$(basename "${svc}")
  mkdir -p "${docDir}"
  $protoCmdPrefix \
  --openapi_out="${docDir}" \
  "${svc}"/*.proto
done

servicePath="${protoSrcBasePath}/spbe/mq"
for svc in "${servicePath}"/*; do
  echo "${svc}"
  dirname=$(basename "${svc}")
  docDir="${outBaseDir}"/spbe/mq/"${dirname}"
  mkdir -p "${docDir}"
  $protoCmdPrefix \
  --doc_out="${docDir}" \
  --doc_opt=html,"${dirname}".html \
  "${svc}"/*.proto
done

servicePath="${protoSrcBasePath}/spbe/common"
for svc in "${servicePath}"/*; do
  echo "${svc}"
  dirname=$(basename "${svc}")
  docDir="${outBaseDir}"/spbe/common/"${dirname}"
  mkdir -p "${docDir}"
  $protoCmdPrefix \
  --openapi_out="${docDir}" \
  --doc_out="${docDir}" \
  --doc_opt=html,"${dirname}".html \
  "${svc}"/*.proto
done

# MOVE TO THEIR DIR
rm -Rfv ./gen/docs
mkdir -p ./gen/docs
mv "${outBaseDir}" ./gen
rm -Rfv "${outBaseDir}"