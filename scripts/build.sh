#/usr/bin/env bash

# get the service we are to build.
SERVICE=$1

# exit if we're not given a service.
if [ -z $SERVICE ]; then
	echo "» no service to build defined, exiting..."
	exit 2
fi

# exit if service doesn't exist.
if [ ! -d "svc/$SERVICE" ]; then
	echo "» service does not exist, exiting..."
	exit 1
fi

# get the current revision from git to embed in the binary.
REVISION=$(git rev-parse HEAD)

# make sure we have a directory to build into.
mkdir -p "svc/$SERVICE/build/dist"

echo "» building binary for $SERVICE at revision $REVISION..."

# build go binary.
go build -o "svc/$SERVICE/build/dist/$SERVICE" -trimpath -ldflags "-X 'gitlab.com/laserdigital/platform/go/pkg/build.revision=$REVISION'" svc/$SERVICE/cmd/main.go

echo "» binary for $SERVICE built to svc/$SERVICE/build/dist/$SERVICE"
