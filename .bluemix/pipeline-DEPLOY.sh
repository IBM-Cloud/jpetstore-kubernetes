#!/bin/bash
INGRESS_HOSTNAME=$(bx cs cluster-get $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressHostname | tr -d '":,' | awk '{print $2}')
echo "INGRESS_HOSTNAME=${INGRESS_HOSTNAME}"

INGRESS_SECRETNAME=$(bx cs cluster-get $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressSecretName | tr -d '":,' | awk '{print $2}')
echo "INGRESS_SECRETNAME=${INGRESS_SECRETNAME}"

# where to put our app
kubectl create namespace $TARGET_NAMESPACE

# copy the tls cert over
kubectl get secret $INGRESS_SECRETNAME -o yaml | sed 's/namespace: default/namespace: '$TARGET_NAMESPACE'/' | kubectl create -f -

# a secret to access the registry
CR_TOKEN_ID=$(ibmcloud cr token-list | grep JPetStore | awk ' {print $1} ')

if [ -z "$CR_TOKEN_ID" ]
then
  CR_TOKEN=$(ibmcloud cr token-add --description "JPetStore toolchain pull token" --non-expiring | sed -n '4,4p' | awk ' {print $2} ')
else
  CR_TOKEN=$(ibmcloud cr token-get $CR_TOKEN_ID | sed -n '4,4p' | awk ' {print $2} ')
fi

kubectl --namespace $TARGET_NAMESPACE create secret docker-registry petstore-docker-registry \
  --docker-server=registry.ng.bluemix.net \
  --docker-password=${CR_TOKEN} \
  --docker-username=token \
  --docker-email=devops@build.com

# create mmssearch secret file
cat > "mms-secrets.json" << EOF
{
  "jpetstoreurl": "http://jpetstore.$INGRESS_HOSTNAME",
  "watson": 
  {
    "url": "https://gateway-a.watsonplatform.net/visual-recognition/api",
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
