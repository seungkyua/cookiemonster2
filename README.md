# cookiemonster2

// we try to use go mod 
make clean
make build-linux
make docker
docker tag cookiemonster seungkyua/cookiemonster
docker push seungkyua/cookiemonster


curl -XPOST http://localhost:8080/api/v1/pod/start
curl -XPOST http://localhost:8080/api/v1/pod/stop

curl -XPOST http://k2-master01:30003/api/v1/pod/start
curl -XPOST http://k2-master01:30003/api/v1/pod/stop
