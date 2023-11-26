# Spendwatch backend 

backend has api endpoints to get expenses from DB; also starts
background thread to update the DB with latest expenses


## To build and run
```
go build -o main *.go
./main
```

## Build docker
```
docker build -t vvksh/spendwatch-docker .
```

## Push

docker push vvksh/spendwatch-docker:latest

## To launch it using docker
 docker run -d -p 8000:8000 vvksh/spendwatch-docker:latest

## Test endpoint on localhost
curl -X GET "localhost:{PORT}/expenses?groupBy=month"