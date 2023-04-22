#!/bin/bash

DEV_ENV="dev"
TEST_ENV="test"
PROD_ENV="prod"
API_ROOT_DOMAIN="testbox.compsoc.ie"
FORCE_CONFIRMED_FILE="/tmp/apid_deployment_manager/forceConfirmed"

# set initial vals
debug=false
delete=false
force=false
environment=""
docker_image=""

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
    echo "$PROD_ENV"
  elif [[ "$1" == *"prerelease"* ]]; then
    echo "$TEST_ENV"
  fi
  echo "$DEV_ENV"
}

function doesContainerExists {
  if ! [ -z "$environment" ]; then 
    echo 'Error: more than one environment set.'
    echo ''
    return 1
  fi
}

function doesImageEnvMatchSelectedEnv {
  if ! [ $(getImageEnv) = "$environment" ] && [ $force = false ]; then 
    echo 'Error: Docker image environment does not match chosen environment.'
    return 1
  fi
  return 0
}

function canForceScript {
  if [ $force = true ] && ! [ -f "$FORCE_CONFIRMED_FILE" ]; then
    touch $FORCE_CONFIRMED_FILE
    echo 'Note: You have choosen to forcefully run this script. If you meant to do this,'
    echo 'please rerun this script. If you didnt mean to do this, delete the forceConfirmed'
    echo "file in $FORCE_CONFIRMED_FILE."
    echo ''
    return 1
  elif [ -f "$FORCE_CONFIRMED_FILE" ]; then
    echo 'Note: You have choosen to forcefully run this script.'
    echo ''
    return 0
  fi
  return 1
}

##############################################
#  VERIFY USER INPUTS
##############################################

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

##############################################
#  
##############################################

if $delete; then
  # delete the container
  exit 0;
fi

#rm $FORCE_CONFIRMED_FILE || true

image_env=$(getImageEnv "$docker_image")
image_tag=$(echo "$docker_image" | cut -d ':' -f 2)
docker_api_domain=""
if [ "prod" = "$image_env" ]; then
  docker_api_domain="$API_ROOT_DOMAIN"
elif [ "test" = "$image_env" ]; then
  docker_api_domain="dev.$API_ROOT_DOMAIN"
else
  docker_api_domain="${image_tag}.dev.$API_ROOT_DOMAIN"
fi

# Check if a container already exists, make a func that can check for existing ones so that delete pocess can use it too
container_name="compsoc_apid_$image_tag"
# Check if container exist and delete
if [ $( docker ps -a -f name=$container_name | wc -l ) -eq 2 ]; then
  docker rm -f $container_name
  echo "Info: Recreating container with name $container_name"
else
  echo "Info: Container with name $container_name doesn't exist. Creating it."
fi

# spin it back up
docker_message=$(docker run --name $container_name --network=web -d $docker_image)
echo $docker_message