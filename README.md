### Simple project example  

Httptest usage for api tests examples: [link](https://github.com/rodkevich/mvpbe/blob/master/internal/domain/sample/controller_test.go)  
Dockertest usage for db tests examples: [link](https://github.com/rodkevich/mvpbe/blob/master/internal/domain/sample/datasource/sample_datasource_test.go)  
  
Docker must be started on your machine to allow dockertest spawn containers.  
You can skip dockertest test-cases with -short flag for tests or through env var settings.  
Example command: ``go test ./... -count=1 -v -short``  
--
#### Scripts:  
to get local env run: `./scripts/local-setup.sh` [link](https://github.com/rodkevich/mvpbe/blob/master/scripts/local-setup.sh)  
to get docker-based env run: `./scripts/docker-setup.sh` [link](https://github.com/rodkevich/mvpbe/blob/master/scripts/docker-setup.sh)  
