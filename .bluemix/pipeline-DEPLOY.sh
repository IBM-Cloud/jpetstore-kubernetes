#!/bin/bash
TARGET_USER=$(ibmcloud target | grep User | awk '{print $2}')
echo "TARGET_USER=$TARGET_USER"

INGRESS_HOSTNAME=$(ibmcloud ks cluster get --cluster $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressHostname | tr -d '":,' | awk '{print $2}')
echo "INGRESS_HOSTNAME=${INGRESS_HOSTNAME}"

INGRESS_SECRETNAME=$(ibmcloud ks cluster get --cluster $PIPELINE_KUBERNETES_CLUSTER_NAME --json | grep ingressSecretName | tr -d '":,' | awk '{print $2}')
echo "INGRESS_SECRETNAME=${INGRESS_SECRETNAME}"

if kubectl get namespace $TARGET_NAMESPACE; then
  echo "Namespace $TARGET_NAMESPACE already exists"
else
  echo "Creating namespace $TARGET_NAMESPACE..."
  kubectl create namespace $TARGET_NAMESPACE || exit 1
fi

# copy the tls cert over
# kubectl get secret $INGRESS_SECRETNAME -o yaml | sed 's/namespace: default/namespace: '$TARGET_NAMESPACE'/' | kubectl create -f -

kubectl get secret all-icr-io -n default -o yaml | sed 's/namespace: default/namespace: '$TARGET_NAMESPACE'/' | kubectl create -f -

# create mmssearch secret file
cat > "mms-secrets.json" << EOF
{
  "jpetstoreurl": "http://jpetstore.$INGRESS_HOSTNAME",
  "watson": 
  {
    "url": "$WATSON_VR_URL",
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

## create mmssearch secret
kubectl delete secret mms-secret --namespace $TARGET_NAMESPACE
kubectl --namespace $TARGET_NAMESPACE create secret generic mms-secret --from-file=mms-secrets=./mms-secrets.json

## Reset tiller
helm reset --force

## Setup tiller 
kubectl create serviceaccount tiller -n kube-system
kubectl create clusterrolebinding tiller --clusterrole=cluster-admin --serviceaccount=kube-system:tiller -n kube-system

## reset and install helm tiller into cluster
helm init --service-account tiller

# install release named jpetstore
helm upgrade --install --namespace $TARGET_NAMESPACE --debug \
  --set image.repository=$REGISTRY_URL/$REGISTRY_NAMESPACE \
  --set image.tag=latest \
  --set image.pullPolicy=Always \
  --set ingress.hosts={jpetstore.$INGRESS_HOSTNAME} \
  --set ingress.secretName=$INGRESS_SECRETNAME \
  --wait jpetstore ./helm/modernpets

# install release named mmssearch
helm upgrade --install --namespace $TARGET_NAMESPACE --debug \
  --set image.repository=$REGISTRY_URL/$REGISTRY_NAMESPACE \
  --set image.tag=latest \
  --set image.pullPolicy=Always \
  --set ingress.hosts={mmssearch.$INGRESS_HOSTNAME} \
  --set ingress.secretName=$INGRESS_SECRETNAME \
  --wait mmssearch ./helm/mmssearch
