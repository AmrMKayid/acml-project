#!/bin/sh

print() {
  printf "\n%b\n" "$1"
}

print "Clone Project From GitHub"
git clone https://github.com/ramin0/myapp

print "Make sure the project is cloned successfully!"
ls

print "Change your current directory to (myapp)"
cd myapp

print "Download Dockerfile from our repo"
curl -OL https://raw.githubusercontent.com/amrmkayid/acml-project/assignment/2/Dockerfile

print "make sure that Dockerfile is created"
ls

print "Build Docker Image"
docker build -t myapp .

print "make sure that docker image is created successfully!"
docker images

print "Run Docker Image"
docker container run myapp