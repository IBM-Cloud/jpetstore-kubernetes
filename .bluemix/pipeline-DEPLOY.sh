#!/bin/bash
TARGET_USER=$(ibmcloud target | grep User | awk '{print $2}')
check_value "$TARGET_USER"
echo "TARGET_USER=$TARGET_USER"

INGRESS_HOSTNAME=$(ibmcloud cs cluster-get $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressHostname | tr -d '":,' | awk '{print $2}')
echo "INGRESS_HOSTNAME=${INGRESS_HOSTNAME}"

INGRESS_SECRETNAME=$(ibmcloud cs cluster-get $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressSecretName | tr -d '":,' | awk '{print $2}')
echo "INGRESS_SECRETNAME=${INGRESS_SECRETNAME}"

if kubectl get namespace $TARGET_NAMESPACE; then
  echo "Namespace $TARGET_NAMESPACE already exists"
else
  echo "Creating namespace $TARGET_NAMESPACE..."
  kubectl create namespace $TARGET_NAMESPACE || exit 1
fi

# copy the tls cert over
# kubectl get secret $INGRESS_SECRETNAME -o yaml | sed 's/namespace: default/namespace: '$TARGET_NAMESPACE'/' | kubectl create -f -

# a secret to access the registry
if kubectl get secret petstore-docker-registry --namespace $TARGET_NAMESPACE; then
  echo "Docker Registry secret already exists"
else
  REGISTRY_TOKEN=$(ibmcloud cr token-add --description "petstore-docker-registry for $TARGET_USER" --non-expiring --quiet)
  kubectl --namespace $TARGET_NAMESPACE create secret docker-registry petstore-docker-registry \
    --docker-server=${REGISTRY_URL} \
    --docker-password="${REGISTRY_TOKEN}" \
    --docker-username=token \
    --docker-email="${TARGET_USER}" || exit 1
fi

# create the imagePullSecret https://cloud.ibm.com/docs/containers/cs_images.html#store_imagePullSecret
kubectl patch -n $TARGET_NAMESPACE serviceaccount/default -p '{"imagePullSecrets":[{"name": "petstore-docker-registry"}]}'

# create mmssearch secret file
cat > "mms-secrets.json" << EOF
{
  "jpetstoreurl": "http://jpetstore.$INGRESS_HOSTNAME",
  "watson": 
  {
    "url": "https://gateway.watsonplatform.net/visual-recognition/api",
    "note": "It may take up to 5 minutes for this key to become active",
    "api_key": "$WATSON_VR_API_KEY"
  },
  "twilio": {
    "sid": "$TWILIO_SID",
    "token": "$TWILIO_TOKEN",
    "number": "$TWILIO_NUMBER"
  }
}
EOF

# create mmssearch secret
kubectl delete secret mms-secret --namespace $TARGET_NAMESPACE
kubectl --namespace $TARGET_NAMESPACE create secret generic mms-secret --from-file=mms-secrets=./mms-secrets.json

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
  --wait jpetstore ./helm/modernpets

# install release named mmssearch
helm upgrade --install --namespace $TARGET_NAMESPACE --debug \
  --set image.repository=$REGISTRY_URL/$REGISTRY_NAMESPACE \
  --set image.tag=latest \
  --set image.pullPolicy=Always \
  --set ingress.hosts={mmssearch.$INGRESS_HOSTNAME} \
  --set ingress.secretName=$INGRESS_SECRETNAME \
  --recreate-pods \
  --wait mmssearch ./helm/mmssearch
