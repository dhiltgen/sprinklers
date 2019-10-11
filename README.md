# Sprinklers

This repo contains a gRPC+REST sprinkler irrigation system aimed at Raspberry PIs
using [relays](https://www.amazon.com/gp/product/B00KTELP3I/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1) hooked up to GPIO pins.

## Hardware

I've tried a few different variations on driving sprinkler circuits over the
years.  My first variant was based on 1-Wire, but that proved to be a bit
complicated and finicky.  For my most recent iteration I switched to a much
simpler model driving relays via GPIO.  The number of circuits I have is less
than the number of available GPIO circuits on the PI, so the added complexity
of MUXing control of multiple circuits was overkill.

<img src="http://www.hiltgen.com/daniel/electronics/sprinkler_controller.jpg" width="400">


## Deploying

The automation in this repo assumes you have a Kubernetes cluster where at
least one node in the cluster has a hostname `sprinklers` which is the target
for the Deployment.

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

