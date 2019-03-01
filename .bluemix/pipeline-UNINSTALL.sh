#!/bin/bash
# delete the two apps
helm delete jpetstore --purge
helm delete mmssearch --purge

# delete the secrets stored in our cluster
kubectl delete secret mms-secret
kubectl delete secret petstore-docker-registry

# remove the images
ibmcloud cr image-rm $DB_IMAGE_URL
ibmcloud cr image-rm $MMS_IMAGE_URL
ibmcloud cr image-rm $WEB_IMAGE_URL

# delete the namespace
kubectl delete namespace $TARGET_NAMESPACE