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
git clone https://github.com/ramin0/myapp 
cd myapp
ls

print "Create & Run Postgres Database Using Docker"
docker run -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=myapp -p 5432:5432 -d postgres:10.1-alpine

sleep 15

print "Execute dump.sql"
cat dump.sql | docker exec -i $( docker ps -q ) psql -U root myapp

print "Build & Run main.go"

print "synchronize external packages using govendor"
govendor sync

print "build main.go using govendor"
govendor build main.go

print "run build file"
./main


print "###############################################"
print "      Separation Methodology: config File      "
print "###############################################"

print "Download config.json example from our repo"
curl -OL https://raw.githubusercontent.com/AmrMKayid/acml-project/assignments/3/config.json-example
mv config.json-example config.json

print "Download new go program from our repo"
curl -OL https://raw.githubusercontent.com/AmrMKayid/acml-project/assignments/3/main_with_config.go

ls

govendor add +e

print "synchronize external packages using govendor"
govendor sync

print "build main.go using govendor"
govendor build main_with_config.go

print "run build file"
./main_with_config



print "###############################################"
print "      Separation Methodology: envs.            "
print "###############################################"

# print "Download config.json example from our repo"
# curl -OL https://raw.githubusercontent.com/AmrMKayid/acml-project/assignments/3/.env-example
# mv .env-example .env

print "Download new go program from our repo"
curl -OL https://raw.githubusercontent.com/AmrMKayid/acml-project/assignments/3/main_with_envs.go

ls

govendor fetch github.com/joho/godotenv
govendor add +e

print "synchronize external packages using govendor"
govendor sync

print "build main.go using govendor"
govendor build main_with_envs.go

print "run build file"
export DATABASE_USERNAME="root" && export DATABASE_PASSWORD="secret" && export DATABASE_NAME="myapp" && ./main_with_envs

