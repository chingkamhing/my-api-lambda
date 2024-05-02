#!/bin/bash
#
# Script file to use curl command to test AWS API Gateway and my-api-lambda.
#

URL="${URL:-"https://on2i8y7rxj.execute-api.ap-southeast-1.amazonaws.com/qa"}"
METHOD="PUT"
ENDPOINT='callback/company-1001'
OPTS="-s --location"
NUM_ARGS=0
DEBUG=""

# Function
SCRIPT_NAME=${0##*/}
Usage () {
	echo
	echo "Description:"
	echo "Script file to use curl command to test AWS API Gateway and my-api-lambda."
	echo
	echo "Usage: $SCRIPT_NAME"
	echo "Options:"
	echo " -d                           Dry run"
	echo " -k                           Allow https insecure connection"
	echo " -u  [url]                    API server URL"
	echo " -p  [port]                   API server port number"
	echo " -v                           Make the curl operation more talkative"
	echo " -h                           This help message"
	echo
}

# Parse input argument(s)
while [ "${1:0:1}" == "-" ]; do
	OPT=${1:1:1}
	case "$OPT" in
	"d")
		DEBUG="echo"
		;;
	"k")
		OPTS="$OPTS -k"
		;;
	"u")
		URL=$2
		shift
		;;
	"p")
		PORT=$2
		shift
		;;
	"v")
		OPTS="$OPTS -v"
		;;
	"h")
		Usage
		exit
		;;
	esac
	shift
done

if [ "$#" -ne "$NUM_ARGS" ]; then
    echo "Invalid parameter!"
	Usage
	exit 1
fi

# trim URL trailing "/"
URL="$(echo -e "${URL}" | sed -e 's/\/*$//')"
if [ "$PORT" == "" ]; then
	url_port="${URL}"
else
	url_port="${URL}:${PORT}"
fi

REQUEST_BODY="$(echo {} | jq \
  --arg hello "Hello, my first API lambda!" \
  '. + { "hello": $hello
       }'
)"

# perform curl
$DEBUG curl $OPTS -d "$REQUEST_BODY" \
	-X $METHOD \
	-H "Content-Type: application/json" \
	-H "Accept: application/json" \
	${url_port}/${ENDPOINT}
