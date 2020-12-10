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
kubectl patch -n ${TARGET_NAMESPACE} serviceaccount/default -p '{"imagePullSecrets":[{"name": "all-icr-io"}]}'

# create mmssearch secret file
cat > "mms-secrets.json" << EOF
{
  "jpetstoreurl": "http://jpetstore.$INGRESS_HOSTNAME",
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
 
## For Helm 3
echo "=========================================================="
echo "CHECKING HELM 3 VERSION (if absent will install latest Helm 3 version) "
set +e
LOCAL_VERSION=$( helm version ${HELM_TLS_OPTION} --template="{{ .Version }}" | cut -c 2- )
# if no Helm 3 locally installed, LOCAL_VERSION will be empty -- will install latest then
set -e
if [ -z "${HELM_VERSION}" ]; then
  CLIENT_VERSION=${LOCAL_VERSION} 
else
  CLIENT_VERSION=${HELM_VERSION}
fi
set +e
if [ -z "${CLIENT_VERSION}" ]; then # Helm 3 not present yet and no explicit required version, install latest
  echo "Installing latest Helm 3 client"
  WORKING_DIR=$(pwd)
  mkdir ~/tmpbin && cd ~/tmpbin
  HELM_INSTALL_DIR=$(pwd)
  curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
  export PATH=${HELM_INSTALL_DIR}:$PATH
  cd $WORKING_DIR
elif [ "${CLIENT_VERSION}" != "${LOCAL_VERSION}" ]; then
  echo -e "Installing Helm 3 client ${CLIENT_VERSION}"
  WORKING_DIR=$(pwd)
  mkdir ~/tmpbin && cd ~/tmpbin
  curl -L https://get.helm.sh/helm-v${CLIENT_VERSION}-linux-amd64.tar.gz -o helm.tar.gz && tar -xzvf helm.tar.gz
  cd linux-amd64
  export PATH=$(pwd):$PATH
  cd $WORKING_DIR
fi
set -e

echo "helm version"
helm version ${HELM_TLS_OPTION}

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