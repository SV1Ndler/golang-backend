#!/usr/bin/bash



SERVERPORT=8082
SERVERADDR=localhost:${SERVERPORT}

FILE_LOC="golang.png"
# FILE_BYTES=`(xxd -b ${FILE_LOC}) | base64`
FILE_BYTES=$(cat "golang.png" | base64 -w0)

set -eux
set -o pipefail
# SMS=$(cat "golang.png" | base64)
#Inject content into JSON
# DATA=${DATA/PLACEHOLDERFORCONTENT/$FILE_BYTES}

# Start by deleting all existing tasks on the server
# curl -iL -w "\n" -X DELETE ${SERVERADDR}/posts/

# Add some tasks
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"title":"title first", "content":"content first",  "created":"2016-01-02T15:04:05+00:00"}' ${SERVERADDR}/posts
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"title":"title second", "content":"content second", "created":"2016-01-02T15:04:05+00:00"}' ${SERVERADDR}/posts
curl -iL -w "\n" -X GET ${SERVERADDR}/posts


curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data "{\"file\":\"${FILE_BYTES}\", \"created\":\"2016-01-02T15:04:05+00:00\"}" ${SERVERADDR}/images
curl -iL -w "\n" -X GET ${SERVERADDR}/images