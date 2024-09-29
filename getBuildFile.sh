#!/bin/bash

cd builds

#build image that creates app file
./build_image.sh

#copy app file to host todo-list/builds
./run_build.sh