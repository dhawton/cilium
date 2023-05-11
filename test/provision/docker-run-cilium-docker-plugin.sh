#!/bin/bash

CILIUM_DOCKER_PLUGIN_IMAGE=${CILIUM_DOCKER_PLUGIN_IMAGE:-cilium/docker-plugin:latest}

set -e
shopt -s extglob

# Run without sudo if not available (e.g., running as root)
SUDO=
if [ ! "$(whoami)" = "root" ] ; then
    SUDO=sudo
fi

if [ "$1" = "uninstall" ] ; then
    if [ -n "$(${SUDO} docker ps -a -q -f label=app=cilium-docker)" ]; then
        echo "Shutting down running Cilium docker plugin"
        ${SUDO} docker rm -f cilium-docker || true
    fi
    if [ -f /usr/bin/cilium-docker ] ; then
        echo "Removing /usr/bin/cilium-docker"
        ${SUDO} rm /usr/bin/cilium-docker
    fi
    exit 0
fi

DOCKER_OPTS+=" --label app=cilium-docker"

if [ -n "$(${SUDO} docker ps -a -q -f label=app=cilium-docker)" ]; then
    echo "Shutting down running Cilium docker-plugin"
    ${SUDO} docker rm -f cilium-docker || true
fi

echo "Launching Cilium docker plugin $CILIUM_DOCKER_PLUGIN_IMAGE with params $CILIUM_OPTS"
${SUDO} docker create --name cilium-docker $DOCKER_OPTS $CILIUM_DOCKER_PLUGIN_IMAGE

# Copy Cilium docker-plugin-binary
${SUDO} docker cp cilium-docker:/usr/bin/cilium-docker /usr/bin/
