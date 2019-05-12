##
WARNING: Work in Progress. 

# LaserProxy
## What is it?
LaserProxy is a HTTP proxy server which understands the configured networks on the system and route intellegently.

Lets say, you have two or more internet connections connected to a machine which is from different ISP providers. Your browser can utilize any one of the internet connection at a given point of time. But using LaserProxy, you can utilize both the connection at the same time. i.e One request can go from one ISP and other request can go from other ISP. 

## Health check
It also does health check via different connections that are connected from the machine and use them which are healthy. 

## Is this Forward Proxy or Reverse Proxy?
LaserProxy can act as both Forward and Reverse Proxies. We can also run more than one proxy server at the same time.

## How to run?
We can start this server using following command. We also run this as init script which will bring the server during system start-up
```bash
sudo ./proxy
```

We need to run this server as root because it is dealing with network devices.

## How to start Forward Proxy?
TODO: Add steps

## How to start Reverse Proxy?
TODO: Add steps


