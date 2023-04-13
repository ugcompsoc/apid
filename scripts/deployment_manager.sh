#!/bin/bash

DEV_ENV="dev"
TEST_ENV="test"
PROD_ENV="prod"
FORCE_CONFIRMED_FILE="/tmp/apid_deployment_manager/forceConfirmed"

# set initial vals
debug=false
delete=false
force=false
environment=""
docker_image=""
api_root_domain="testbox.compsoc.ie"
docker_image_env=""

function info {
  echo 'University Of Galway Computer Society API Deployer'
  echo '  '
  echo '  This bash script is used to deploy various versions of our API, including:'
  echo '    - development versions (git branches)'
  echo '    - unstable prereleases (tags of main)'
  echo '    - stable releases (git releases)'
  echo '  '
  echo '  For the most part, this script takes the specified images from DockerHub and'
  echo '  deploys them as a docker container. Traefik then takes over and will route'
  echo '  requests to the right environments. E.g. if you want to access your'
  echo '  development branch you would access them at {branch}.dev.api.{domain}, tags'
  echo '  are promoted to dev.api.{domain}. Stable releases are pushed to api.{domain}.'
  echo '  '
  echo '  Each container is also given their own supporting containers, such as there own'
  echo '  MongoDB. The aim will be that these container will be a copy of what is in prod'
  echo '  but will have the PII (personally identifiable information) removed.'
  echo '  '
  echo '  Enjoy yourself - Sláinte.'
  echo '  '
  exit 0
}

function usage {
  echo "Usage: $(basename $0) [-BTP] --image IMAGE_TAG [--delete] [--debug] [--help]" 2>&1
  echo '  -B | --branch-env    deploy to branch/development environment'
  echo '  -T | --test-env      deploy to test environment'
  echo '  -P | --prod-env      deploy to prod environment'
  echo '  -I | --image         dockerhub image tag'
  echo '  -D | --delete        delete existing container'
  echo '  -d | --debug         print debug messages'
  echo '  -h | --help          print this'
  echo '  '
  echo 'Example:'
  echo "  $(basename $0) -B --image ugcompsoc/apid:latest"
  echo '  '
  exit 1
}

function exitIfEnvSet {
  if ! [ -z "$environment" ]; then 
    echo 'Error: more than one environment set.'
    echo ''
    exit 1
  fi
}

function getImageEnv() {
  if [[ "$1" =~ "^v[0-9]$" ]]; then
    docker_api_image_env=$PROD_ENV
  elif [[ "$1" == *"prerelease"* ]]; then
    docker_api_image_env=$TEST_ENV
  else
    docker_api_image_env=$DEV_ENV
  fi
}

function containerExists {
  if ! [ -z "$environment" ]; then 
    echo 'Error: more than one environment set.'
    echo ''
    exit 1
  fi
}

if [[ ${#} -eq 0 ]]; then
  info
fi

# Option strings
SHORT=BTPi:Dfdh
LONG=branch-env,test-env,prod-env,image:,delete,force,debug,help

# read the options
OPTS=$(getopt --options $SHORT --long $LONG -- "$@")

eval set -- "$OPTS"

while true ; do
  case "$1" in
    -B | --branch-env )
      exitIfEnvSet
      environment="$DEV_ENV"
      shift
      ;;
    -T | --test-env )
      exitIfEnvSet
      environment="$TEST_ENV"
      shift
      ;;
    -P | --prod-env )
      exitIfEnvSet
      environment="$PROD_ENV"
      shift
      ;;
    -i | --image )
      docker_image="$2"
      shift 2
      ;;
    -D | --delete )
      delete=true
      shift
      ;;
    -f | --force )
      force=true
      shift
      ;;
    -d | --debug )
      debug=true
      shift
      ;;
    -h | --help )
      usage
      shift
      ;;
    -- )
      shift
      break
      ;;
    *)
      echo "Internal error!"
      echo ''
      exit 1
      ;;
  esac
done

if [ -z "$environment" ]; then 
  echo 'Error: Environment is not set.'
  echo ''
  exit 1
fi

if [ -z "$docker_image" ]; then 
  echo 'Error: Docker image not set.'
  echo ''
  exit 1
fi

if $delete; then
  # delete the container
  exit 0;
fi

if ! [ "$docker_api_image_env" = "$environment" ] && [ $force = false ]; then 
  echo 'Error: Docker image environment does not match chosen environment.'
  echo ''
  exit 1
elif [ $force = true ] && ! [ -f "$FORCE_CONFIRMED_FILE" ]; then
  touch /tmp/forceConfirmed
  echo 'Note: You have choosen to forcefully run this script. If you meant to do this,'
  echo 'please rerun this script. If you didnt mean to do this, delete the forceConfirmed'
  echo 'file.'
  echo ''
  exit 0
fi

rm $FORCE_CONFIRMED_FILE || true

getImageEnv "$docker_api_image_tag"
docker_api_domain=""
if [ "prod" = "$docker_api_image_env" ]; then
  docker_api_domain="${api_root_domain}";
elif [ "test" = "$docker_api_image_env"  ]; then
  docker_api_domain="dev.${api_root_domain}";
else
  docker_api_domain="${docker_api_image_tag}.dev.${api_root_domain}";
fi

# Check if a container already exists, make a func that can check for existing ones so that delete pocess can use it too
# Delete it and it's db or whatever else
# spin it back up
