# cookiemonster2

brew install glide
glide create

make build
make docker

curl -XPOST http://localhost:8080/api/v1/pod/start
curl -XPOST http://localhost:8080/api/v1/pod/stop