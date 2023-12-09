# Example go http project
> For create go http server.

### Binary build
```bash
CGO_ENABLED=0 GOOS=linux go build -o worker

CGO_ENABLED=0 GOOS=linux go build -o backend

CGO_ENABLED=0 GOOS=linux go build -o frontend
``` 

### Container Image build
```bash
docker build -f Dockerfile.worker -t ghcr.io/neilkuan/example-go-http:worker .

docker build -f Dockerfile.frontend -t ghcr.io/neilkuan/example-go-http:frontend .

docker build -f Dockerfile.backend -t ghcr.io/neilkuan/example-go-http:backend .
``` 


### Deploy to kubernetes from ghcr.io
```bash
kubectl create deploy worker --image ghcr.io/neilkuan/example-go-http:worker --port=8080
kubectl create deploy frontend --image ghcr.io/neilkuan/example-go-http:frontend --port=8080
kubectl create deploy backend --image ghcr.io/neilkuan/example-go-http:backend --port=8080
```
