#!/bin/sh

print() {
  printf "\n%b\n" "$1"
}

print "Install go & govendor"
apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go
wget -O go.tgz https://dl.google.com/go/go1.10.3.src.tar.gz
tar -C /usr/local -xzf go.tgz
cd /usr/local/go/src/
./make.bash
export PATH="/usr/local/go/bin:$PATH" 
export GOPATH=/opt/go/
export PATH=$PATH:$GOPATH/bin
apk del .build-deps

print "Make sure it is installed successfully:"
go version

print "Install govendor:"
go get -u github.com/kardianos/govendor

print "Clone Project From Github"
git clone https://github.com/ramin0/myapp && cd myapp && ls

print "Create & Run Postgres Database Using Docker"
docker run -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=myapp -p 5432:5432 -d postgres:10.1-alpine

print "Get running container(s) ID"
$container_id = $(docker ps -a -q)

print "Execute dump.sql"
cat dump.sql | docker exec -i $container_id psql -U root myapp

print "Build & Run main.go"

print "synchronize external packages using govendor"
govendor sync

print "build main.go using govendor"
govendor build main.go

print "run build file"
./main









