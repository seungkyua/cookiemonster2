# Cookiemonster
 
 
 * This program delete pod randomly in all namespace
 * It will helpful for test your k8s cluster stability
 * This tool also support feature that shutdown your cluster node by redfish, but it demands more developing
 * Now I develop operator feature
 * Helm chart will be developed

## For install
```sh 
go mod vendor  
make build-linux  
make docker 
```

## Docker hub update 
```sh
docker tag cookiemonster seungkyua/cookiemonster  
docker push seungkyua/cookiemonster  
```

## Start in local environment  
```sh
go run server.go  
curl -X POST http://localhost:10080/api/v1/pod/start  
curl -X POST http://localhost:10080/api/v1/pod/stop  
```

## Start in k8s cluster
```sh
cd cookiemonster-operator && kubectl apply -f deploy
cd deploy && kubectl apply -f crds
curl -X POST http://localhost:${server_port_num}/api/v1/pod/start
curl -X POST http://localhost:${server_port_num}/api/v1/pod/stop
```
## Restart nodes
```sh
modify bmc value in cr file in crd folder 
curl -X POST http://localhost:${server_port_num}/api/v1/node/start
It has to be developed(Stop feature)
```
