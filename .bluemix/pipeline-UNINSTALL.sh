#!/bin/bash
# delete the two apps
helm delete jpetstore --purge
helm delete mmssearch --purge

# delete the secrets stored in our cluster
kubectl delete secret mms-secret
