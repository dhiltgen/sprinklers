# sprinklers
Simple REST app for raspberry pi GPIO based sprinklers

A few build variants...
```
docker -c sprinklers build -t dhiltgen/sprinklers .
```
```
docker -c pi64 buildx build --platform linux/arm/v7 -t dhiltgen/sprinklers:armv7 --push .
```

Recreate the sprinklerd server
```
docker -c sprinklers rm -f sprinklers-grpc
docker -c sprinklers run --rm -d --restart=always --name sprinklers-grpc --privileged -p 1600:1600 -p 1601:1601 dhiltgen/sprinklers:armv7
```

Start a circuit with a timeout
```
docker -c sprinklers run --rm --entrypoint sprinklers dhiltgen/sprinklers:armv7 update --start --stop-after 5m "lawn front"
```


On proto changes for local dev...
```
(cd api && protoc ./sprinklers.proto --go_out=plugins=grpc:sprinklers)
```
