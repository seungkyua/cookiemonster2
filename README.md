# cookiemonster2

brew install glide
glide create

go build -o bin/cookiemonster-linux-amd64 -v ./src/cmd/server.go
make docker
docker tag cookiemonster seungkyua/cookiemonster
docker push seungkyua/cookiemonster


curl -XPOST http://localhost:8080/api/v1/pod/start
curl -XPOST http://localhost:8080/api/v1/pod/stop

curl -XPOST http://k2-master01:30003/api/v1/pod/start
curl -XPOST http://k2-master01:30003/api/v1/pod/stop