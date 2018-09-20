#!/bin/bash
source $ARCHIVE_DIR/build.properties

INGRESS_HOSTNAME=$(bx cs cluster-get $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressHostname | tr -d '":,' | awk '{print $2}')
echo "INGRESS_HOSTNAME=${INGRESS_HOSTNAME}"

INGRESS_SECRETNAME=$(bx cs cluster-get $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressSecretName | tr -d '":,' | awk '{print $2}')
echo "INGRESS_SECRETNAME=${INGRESS_SECRETNAME}"

# where to put our app
kubectl create namespace $TARGET_NAMESPACE

# copy the tls cert over
kubectl get secret $INGRESS_SECRETNAME -o yaml | sed 's/namespace: default/namespace: '$TARGET_NAMESPACE'/' | kubectl create -f -

# a secret to access the registry
kubectl --namespace $TARGET_NAMESPACE create secret docker-registry petstore-docker-registry \
  --docker-server=registry.ng.bluemix.net \
  --docker-password=${PIPELINE_BLUEMIX_API_KEY} \
  --docker-username=iamapikey \
  --docker-email=devops@build.com

## install helm tiller into cluster
helm init

# install release named jpetstore
helm upgrade --install --namespace $TARGET_NAMESPACE --debug \
  --set image.repository=$REGISTRY_URL/$REGISTRY_NAMESPACE \
  --set image.tag=latest \
  --set image.pullPolicy=Always \
  --set ingress.hosts={jpetstore.$INGRESS_HOSTNAME} \
  --set ingress.secretName=$INGRESS_SECRETNAME \
  --recreate-pods \
  --wait jpetstore ../helm/modernpets

# install release named mmssearch
helm upgrade --install --namespace $TARGET_NAMESPACE --debug \
  --set image.repository=$REGISTRY_URL/$REGISTRY_NAMESPACE \
  --set image.tag=latest \
  --set image.pullPolicy=Always \
  --set ingress.hosts={mmssearch.$INGRESS_HOSTNAME} \
  --set ingress.secretName=$INGRESS_SECRETNAME \
  --recreate-pods \
  --wait mmssearch ../helm/mmssearch
