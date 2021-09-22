#!/bin/bash
# fail the script if a command fails
set -eo pipefail

echo -e "Build environment variables:"
echo "REGISTRY_URL=${REGISTRY_URL}"
echo "REGISTRY_NAMESPACE=${REGISTRY_NAMESPACE}"
echo "BUILD_NUMBER=${BUILD_NUMBER}"
echo "ARCHIVE_DIR=${ARCHIVE_DIR}"

# Checking ig buildctl is installed
if which buildctl > /dev/null 2>&1; then
  buildctl --version
else 
  echo "Installing Buildkit builctl"
  curl -sL https://github.com/moby/buildkit/releases/download/v0.8.1/buildkit-v0.8.1.linux-amd64.tar.gz | tar -C /tmp -xz bin/buildctl && mv /tmp/bin/buildctl /usr/bin/buildctl && rmdir --ignore-fail-on-non-empty /tmp/bin
  buildctl --version
fi

WEB_IMAGE_URL=$REGISTRY_URL/$REGISTRY_NAMESPACE/$WEB_IMAGE_NAME
DB_IMAGE_URL=$REGISTRY_URL/$REGISTRY_NAMESPACE/$DB_IMAGE_NAME
MMS_IMAGE_URL=$REGISTRY_URL/$REGISTRY_NAMESPACE/$MMS_IMAGE_NAME

# Create the config.json file to make private container registry accessible
export DOCKER_CONFIG=$(mktemp -d -t cr-config-XXXXXXXXXX)
kubectl create secret --dry-run=true --output=json \
  docker-registry registry-dockerconfig-secret \
  --docker-server=${REGISTRY_URL} \
  --docker-password=${PIPELINE_BLUEMIX_API_KEY} \
  --docker-username=iamapikey --docker-email=a@b.com | \
jq -r '.data[".dockerconfigjson"]' | base64 -d > ${DOCKER_CONFIG}/config.json

echo "=========================================================="
echo -e "BUILDING CONTAINER IMAGE: ${WEB_IMAGE_NAME}"
set -x
buildctl build \
    --frontend dockerfile.v0 --opt filename=Dockerfile --local dockerfile=./jpetstore \
    --local context=./jpetstore \
    --import-cache type=registry,ref=${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${WEB_IMAGE_NAME} \
    --output type=image,name="${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${WEB_IMAGE_NAME}",push=true
set +x
ibmcloud cr image-inspect ${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${WEB_IMAGE_NAME}
ibmcloud cr images --restrict ${REGISTRY_NAMESPACE}/${WEB_IMAGE_NAME}

echo "=========================================================="
echo -e "BUILDING CONTAINER IMAGE: ${DB_IMAGE_NAME}"
set -x
buildctl build \
    --frontend dockerfile.v0 --opt filename=Dockerfile --local dockerfile=./jpetstore/db \
    --local context=./jpetstore/db \
    --import-cache type=registry,ref=${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${DB_IMAGE_NAME} \
    --output type=image,name="${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${DB_IMAGE_NAME}",push=true
set +x
ibmcloud cr image-inspect ${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${DB_IMAGE_NAME}
ibmcloud cr images --restrict ${REGISTRY_NAMESPACE}/${DB_IMAGE_NAME}

echo "=========================================================="
echo -e "BUILDING CONTAINER IMAGE: ${MMS_IMAGE_NAME}"
set -x
buildctl build \
    --frontend dockerfile.v0 --opt filename=Dockerfile --local dockerfile=./mmssearch \
    --local context=./mmssearch \
    --import-cache type=registry,ref=${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${MMS_IMAGE_NAME} \
    --output type=image,name="${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${MMS_IMAGE_NAME}",push=true
set +x
ibmcloud cr image-inspect ${REGISTRY_URL}/${REGISTRY_NAMESPACE}/${MMS_IMAGE_NAME}
ibmcloud cr images --restrict ${REGISTRY_NAMESPACE}/${MMS_IMAGE_NAME}

mkdir -p $ARCHIVE_DIR
echo "REGISTRY_URL=${REGISTRY_URL}" >> $ARCHIVE_DIR/build.properties
echo "REGISTRY_NAMESPACE=${REGISTRY_NAMESPACE}" >> $ARCHIVE_DIR/build.properties
echo "BUILD_NUMBER=${BUILD_NUMBER}" >> $ARCHIVE_DIR/build.properties
echo "DB_IMAGE_URL=${DB_IMAGE_URL}" >> $ARCHIVE_DIR/build.properties
echo "MMS_IMAGE_URL=${MMS_IMAGE_URL}" >> $ARCHIVE_DIR/build.properties
echo "WEB_IMAGE_URL=${WEB_IMAGE_URL}" >> $ARCHIVE_DIR/build.properties
