# Sprinklers

This repo contains a gRPC+REST sprinkler system aimed at Raspberry PIs
using [relays](https://www.amazon.com/gp/product/B00KTELP3I/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1) hooked up to GPIO pins.

## Hardware

I've tried a few different variations on driving sprinkler circuits.  My most
recent incarnation and what I think is the simplest to get set up and running
is based on a very simple GPIO based relay bank.


## Deploying

The automation in this repo assumes a Kubernetes cluster where at least one
node in the cluster has a hostname "sprinklers" which is the target for the
Deployment.

```sh
make deploy
```

Interacting with the sprinklers
```sh
kubectl run -i -t --rm --image=dhiltgen/sprinklers --restart=Never -- sprinklers --help
```
```sh
% kubectl run -i -t --rm --image=dhiltgen/sprinklers --restart=Never -- sprinklers list
If you don't see a command prompt, try pressing enter.
NAME    DESCRIPTION       WATERING NOW    TIME REMAINING
6       Circuit 4 Left    false           
5       Circuit 5 Left    false           
21      Circuit 6 Left    false           
20      Circuit 7 Left    false           
16      Circuit 8 Left    false           
7       blueberries       false           
8       lawn middle       true            seconds:180 
25      roses             false           
24      back fence        false           
22      lawn front        false           
10      garden            false           
9       lawn back         false           
pod "sprinklers" deleted
```

```sh
% kubectl run -i -t --rm --image=dhiltgen/sprinklers --restart=Never -- sprinklers update --start --stop-after 3m "lawn middle"
If you don't see a command prompt, try pressing enter.
NAME    DESCRIPTION    WATERING NOW    TIME REMAINING
8       lawn middle    true            seconds:180 
pod "sprinklers" deleted
```

## Development Setup

This repo leverages [buildx](https://github.com/docker/buildx) to easily build
multi-architecture Docker images.  Assuming your main development environment
is x86, before you can take advantage of buildx for multi-arch builds, you'll
need to add contexts to at least one remote Raspberry PI.  As long as you're
running a recent engine, the `ssh` remote access mechanism is the easiest to
set up (you have to set up your `authorized_keys` so there's no password
prompt):

```sh
docker context create raspberrypi --docker "host=ssh://pi@raspberrypi"
```

**WARNING** Avoid using the same Raspberry PI for building and running in your
Kubernetes cluster.  Kube is designed to run without swap and keep tight
constraints on memory usage, but building can chew through a ton of RAM and
swap helps avoid failed builds.

With that you can create a multi-arch builder with something like the following:
```sh
docker buildx create --use --name mybuild local
docker buildx create --append --name mybuild raspberrypi
```

Then you can build: (**WARNING** - this will build and push to Docker Hub)
```
make ORG=yourhubid
```

