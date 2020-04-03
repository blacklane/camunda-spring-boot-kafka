#!/usr/bin/env bash

set -e

if [[ -z $1 ]]; then
  echo -e  "\nUsage: \n   bin/destroy-custom-stack.sh NAME-OF-YOUR-CUSTOM-STACK\n"
  exit 1
fi

case $1 in
  production|testing|sandbox)
    echo "${1} is a forbidden environment. Aborting."
    exit 1
    ;;
esac

if ! [[ $(aws sts get-caller-identity --profile testing 2>/dev/null) ]]; then
  echo "Your aws credentials are not valid, please renew them!"
  exit 1
fi

command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is required; Run `blops setup` to install it. Aborting."; exit 1; }

NAMESPACE=$1

AWS_PROFILE=testing kubectl --context=testing delete --ignore-not-found deployment,configmap,secret,ingress,svc -n $NAMESPACE chat-sync chat-sync-db
AWS_PROFILE=testing kubectl --context=testing delete --ignore-not-found pod -n $NAMESPACE chat-sync-migrate

# enable this to remove your custom envs for the namespace
# rm -r deploy/overlays/$NAMESPACE

# sed -i doesn't work the same on unix and linux
sys=$(uname -s)
if [[ $sys == *"arwin" ]]; then
  sed -i '' '/^'${NAMESPACE}':/,$d' deploy/deployments.yaml
else
  sed -i '/^'${NAMESPACE}':/,$d' deploy/deployments.yaml
fi
