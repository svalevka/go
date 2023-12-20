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

# build the binary.
./scripts/build.sh $SERVICE

echo "» building package for $SERVICE..."

# build the package.
cd svc/$SERVICE
nfpm package --packager rpm --config nfpm.yml --target build/dist/

echo "» package for $SERVICE built to svc/$SERVICE/build/dist"
