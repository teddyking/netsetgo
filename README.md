# netsetgo

On your marks, net set, GO!

## About

netsetgo is a small binary that helps to setup network namespaces for containers.
It achieves this by:

* Creating a bridge device in the host's network namespace
* Creating a veth pair
  * One side of the pair is attached to the bridge
  * The other side is placed inside the container's network namespace
* Ensuring any traffic originating from the container's network namespace is routed via the veth

Note: netsetgo is intended for demonstration/learning purposes only, and not as a production-ready container networker.

## Usage

```
sudo netsetgo \
  -bridgeAddress 10.10.10.1/24 \
  -bridgeName brg0 \
  -containerAddress 10.10.10.2/24 \
  -pid 100 \
  -vethNamePrefix veth
```
