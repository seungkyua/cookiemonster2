# cookiemonster2
 
 
 * This program delete pod randomly
 * It will be update node version 
 
//For install  
go mod vendor  
make clean  
make build-linux  
make docker  

//Docker hub update  
docker tag cookiemonster seungkyua/cookiemonster  
docker push seungkyua/cookiemonster  

//Start in local environment  
go run server.go  
curl -XPOST http://localhost:10080/api/v1/pod/start  
curl -XPOST http://localhost:10080/api/v1/pod/stop  

