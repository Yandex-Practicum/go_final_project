#!/bin/bash

docker run --name todo_app_build --rm -d todo_app:build sleep 180

docker cp todo_app_build:/app/builds/todo_app ./todo_app

docker stop todo_app_build