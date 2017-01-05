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

The `-pid` argument is used to determine the Network namespace in which to place the container's side of the veth pair.
All other args are optional.

## Limitations

netsetgo does not currently perform any iptables configuration, which could mean
that your containers aren't able to reach the Internet. The following set of
iptables rules will enable Internet connectivity.

```
iptables -tnat -N netsetgo
iptables -tnat -A PREROUTING -m addrtype --dst-type LOCAL -j netsetgo
iptables -tnat -A OUTPUT ! -d 127.0.0.0/8 -m addrtype --dst-type LOCAL -j netsetgo
iptables -tnat -A POSTROUTING -s 10.10.10.0/24 ! -o brg0 -j MASQUERADE
iptables -tnat -A netsetgo -i brg0 -j RETURN
```

Also note that netsetgo does not currently perform any cleanup of devices. Veth devices should
be removed automatically once they are no longer in use, but bridge devices will probably hang around.
Finally, netsetgo isn't built to be run concurrently.

## Testing

[concourse CI](http://concourse.ci/) is used for both testing and CI.
The test suite can be run against a local concourse-lite machine with the following:

```
fly -t lite e -c ci/test.yml -x -p -i netsetgo-src=.
```

## Special Thanks

netsetgo would not be possible without the
[netlink](https://github.com/vishvananda/netlink) and
[kawasaki](https://github.com/cloudfoundry/guardian/tree/master/kawasaki)
packages.
