#!/usr/bin/env bash

set -e

SKIP_OVERRIDE=0
NAMESPACE=
OPTSPEC=":n:sh-:"
while getopts "$OPTSPEC" optchar; do
    case "${optchar}" in
        -)
            case "${OPTARG}" in
                skip-override)
                    SKIP_OVERRIDE=1
                    ;;
                namespace=*)
										NAMESPACE=${OPTARG#*=}
                    ;;
                *)
                    if [ "$OPTERR" = 1 ] && [ "${OPTSPEC:0:1}" != ":" ]; then
                        echo "Unknown option --${OPTARG}" >&2
                    fi
                    ;;
            esac;;
        h)
            echo "usage: $0 -n <my-namespace> [--namespace=<my-namespace>] [-s] [--skip-override]" >&2
            exit 2
            ;;
        n)
            NAMESPACE=$OPTARG
            ;;
        s)
        		SKIP_OVERRIDE=1
            ;;
        *)
            if [ "$OPTERR" != 1 ] || [ "${OPTSPEC:0:1}" = ":" ]; then
                echo "Non-option argument: '-${OPTARG}'" >&2
            fi
            ;;
    esac
done

if [[ -z $NAMESPACE ]]; then
  echo -e  "\nUsage: \n   $0 -n <my-namespace> [--namespace=<my-namespace>] [-s] [--skip-override]\n"
  exit 1
fi

case $NAMESPACE in
  production|testing|sandbox)
    echo "${1} is a forbidden environment. Aborting."
    exit 1
    ;;
esac

if ! [[ $(aws sts get-caller-identity --profile testing 2>/dev/null) ]]; then
  echo "Your aws credentials are not valid, please renew them!"
  exit 1
fi

command -v jq >/dev/null 2>&1 || { echo >&2 "jq is required; install it and try again. Aborting."; exit 1; }
command -v yq >/dev/null 2>&1 || { echo >&2 "yq is required; install it and try again (e.g. `pip install --user yq`). Aborting."; exit 1; }
command -v blops >/dev/null 2>&1 || { echo >&2 "blops is required; Check out https://github.com/blacklane/blops and install it. Aborting."; exit 1; }

# create the directory only if it's not present yet. Otherwise just update the containing files
if [ ! -d deploy/overlays/$NAMESPACE ]; then
  cp -r deploy/overlays/custom-stack/ deploy/overlays/$NAMESPACE
elif [ $SKIP_OVERRIDE -eq 0 ]; then
	echo "Updating your custom stack env overlays..."
  cp -r deploy/overlays/custom-stack/* deploy/overlays/$NAMESPACE
fi

# sed -i doesn't work the same on unix and linux
sys=$(uname -s)
if [[ $sys == *"arwin" ]]; then
  find deploy/overlays/$NAMESPACE -type f -exec sed -i '' $string 's/custom-stack/'${NAMESPACE}'/g' {} \;
else
  find deploy/overlays/$NAMESPACE -type f -exec sed -i $string 's/custom-stack/'${NAMESPACE}'/g' {} \;
fi

# copy deployment steps for custom stack namespace if required
custom_stack_config=$(cat deploy/deployments.yaml | grep ${NAMESPACE})
if [[ $custom_stack_config = *[![:space:]]* ]]; then
	echo "Skipping stack config creation..."
else
	# copy the custom-stack config, modify and append to file
	echo "Copying custom stack config..."
	config=$(cat deploy/deployments.yaml \
	 | yq -y .\"custom-stack\" \
	 | sed -e "s/custom-stack/${NAMESPACE}/" \
	 | sed -e "s/^/  /" \
	 | sed -e "s/-/  -/")

	echo -e ${NAMESPACE}":\n$config\n" >> deploy/deployments.yaml
fi

sha=$(git rev-parse --short HEAD)
branch=$(git rev-parse --abbrev-ref HEAD | sed 's/_/-/g' | tr "[:upper:]" "[:lower:]")

if [ "$branch" = "master" ]; then
  docker_tag="master";
else
  docker_tag="pr-$branch-$sha";
fi

echo -e "Pulling image with tag \e[96m${docker_tag}"

AWS_PROFILE=testing blops deploy ${NAMESPACE} -t $docker_tag
if [ $? -eq 0 ]; then
	echo -e "\n# Done! See service in k8s with: \nAWS_PROFILE=testing kubectl --context=testing get pods -n ${NAMESPACE}"
	echo -e "# Enjoy testing!!"
else
	echo "Failed to execute deployment"
	exit 1
fi
