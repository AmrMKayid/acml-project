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

print "Clone repository"
git clone https://github.com/ramin0/myapp

cd myapp

ls

rm docker-compose.yml
rm main.go

curl -OL https://raw.githubusercontent.com/AmrMKayid/acml-project/assignments/4/docker-compose.yml
curl -OL https://raw.githubusercontent.com/AmrMKayid/acml-project/assignments/4/main_with_redis.go
mv main_with_redis.go main.go

govendor fetch github.com/go-redis/redis
govendor add +e

print "synchronize external packages using govendor"
govendor sync

docker-compose up