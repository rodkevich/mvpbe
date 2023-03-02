### Simple project example  

Httptest usage for api tests examples: [link](https://github.com/rodkevich/mvpbe/blob/master/internal/domain/item/controller_test.go)  
Dockertest usage for db tests examples: [link](https://github.com/rodkevich/mvpbe/blob/master/internal/domain/item/datasource/sample_datasource_test.go)  

### Docker:  
Docker must be started on your machine to allow dockertest spawn containers.  
You can skip dockertest test-cases with -short flag for tests or through env var settings.  
Example command: 

    go test ./... -count=1 -v -short

### Scripts:
to get local env run: `./scripts/local-setup.sh` [link](https://github.com/rodkevich/mvpbe/blob/master/scripts/local-setup.sh)  
to get docker-based env run: `./scripts/docker-setup.sh` [link](https://github.com/rodkevich/mvpbe/blob/master/scripts/docker-setup.sh)

### Load tests:  
You can run some load tests using wrk  

To create items in db use POST method like:  
    
    wrk "http://localhost:8080/api/v1/items" -s ./examples/post_binary.lua --latency -t 5 -c 10 -d 10  

To update items in db use PUT method like:  
    
    wrk "http://localhost:8080/api/v1/items/1" -s ./examples/put_binary.lua --latency -t 5 -c 10 -d 10``

#### Else you can use CURL like:
    
    curl -X PUT --location "http://localhost:8080/api/v1/items/450" \
    -H "Content-Type: application/json" \
    -d "{\"status\": \"CREATED\"}"