#!/bin/bash

WORKSPACE="$(pwd)"
CRED_JSON_FILE="$HOME/workspace/jinko/jinko-devops-c0a195928a92.json"

docker run -it --rm --pull=always \
    -e SANDBOX_USER_ID=$(id -u) \
    -e WORKSPACE_MOUNT_PATH=$WORKSPACE \
    -e SANDBOX_RUNTIME_CONTAINER_IMAGE=docker.all-hands.dev/all-hands-ai/runtime:0.39-nikolaik \
    -e LLM_NUM_RETRIES=100 \
    -e LLM_RETRY_MIN_WAIT=60 \
    -e LLM_RETRY_MAX_WAIT=60 \
    -e LLM_RETRY_MULTIPLIER=1 \
    -e LOG_ALL_EVENTS=true \
    -v $WORKSPACE:/opt/workspace_base \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.openhands-state:/.openhands-state \
    -p 13000:3000 \
    --add-host host.docker.internal:host-gateway \
    --name openhands-affiliate-backend \
    docker.all-hands.dev/all-hands-ai/openhands:0.39
