workpath=`pwd`
cd ./src/services/$1
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o 'bin/server' -mod=vendor ./cmds/service

if [ -x /tmp/.insomnia/$1 ]; then
  rm -rf /tmp/.insomnia/$1
fi
mkdir -p /tmp/.insomnia/$1
mkdir -p /tmp/.insomnia/$1/bin
mkdir -p /tmp/.insomnia/$1/config
mkdir -p /tmp/.insomnia/$1/proto

cd $workpath
cp ./src/services/$1/bin/server /tmp/.insomnia/$1/bin/
cp ./src/services/$1/*.yaml /tmp/.insomnia/$1/
cp ./src/services/$1/config/*.yaml /tmp/.insomnia/$1/config/
cp ./src/services/$1/proto/*.proto /tmp/.insomnia/$1/proto/
cp scripts/docker/entrypoint.sh /tmp/.insomnia/$1/

docker build -t $2 \
--build-arg APPID=$1 \
--no-cache \
-f ./scripts/docker/Dockerfile /tmp/.insomnia

# clean
if [ -x /tmp/.insomnia/$1 ]; then
  rm -rf /tmp/.insomnia/$1
fi

# sh scripts/docker/build.sh   rule   insomnia-rule:test
#                              appid     tag_name