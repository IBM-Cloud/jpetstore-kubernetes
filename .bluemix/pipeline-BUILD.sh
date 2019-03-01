#!/bin/bash
echo -e "Build environment variables:"
echo "REGISTRY_URL=${REGISTRY_URL}"
echo "REGISTRY_NAMESPACE=${REGISTRY_NAMESPACE}"
echo "BUILD_NUMBER=${BUILD_NUMBER}"
echo "ARCHIVE_DIR=${ARCHIVE_DIR}"

DB_IMAGE_URL=$REGISTRY_URL/$REGISTRY_NAMESPACE/$DB_IMAGE_NAME
MMS_IMAGE_URL=$REGISTRY_URL/$REGISTRY_NAMESPACE/$MMS_IMAGE_NAME
WEB_IMAGE_URL=$REGISTRY_URL/$REGISTRY_NAMESPACE/$WEB_IMAGE_NAME

set -e
ibmcloud cr build ./jpetstore -t $WEB_IMAGE_URL
ibmcloud cr build ./jpetstore/db -t $DB_IMAGE_URL
ibmcloud cr build ./mmssearch -t $MMS_IMAGE_URL
set +e

mkdir -p $ARCHIVE_DIR
echo "REGISTRY_URL=${REGISTRY_URL}" >> $ARCHIVE_DIR/build.properties
echo "REGISTRY_NAMESPACE=${REGISTRY_NAMESPACE}" >> $ARCHIVE_DIR/build.properties
echo "BUILD_NUMBER=${BUILD_NUMBER}" >> $ARCHIVE_DIR/build.properties
echo "DB_IMAGE_URL=${DB_IMAGE_URL}" >> $ARCHIVE_DIR/build.properties
echo "MMS_IMAGE_URL=${MMS_IMAGE_URL}" >> $ARCHIVE_DIR/build.properties
echo "WEB_IMAGE_URL=${WEB_IMAGE_URL}" >> $ARCHIVE_DIR/build.properties
