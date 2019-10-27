# HTTP Relay
Client to client HTTP communication

**IMPORTANT!!!**
_Httprelay is moving from SaS to OSS.
Service that is accessible at https://httprelay.io is moved to https://demo.httprelay.io.
Please use your own deployments as https://demo.httprelay.io for demonstration purpose only!_ 

## Features
- Communication between two or more HTTP clients (e.g. web browsers, mobile apps, IoT devices etc.)
- Peer to peer mode.
- Multicast (many receivers) mode.

## Installation
### Download
[Download executable](https://gitlab.com/jonas.jasas/httprelay/-/jobs/artifacts/master/browse/download?job=build:download) (Linux, Mac, Windows)

### Docker
- Latest image: `registry.gitlab.com/jonas.jasas/httprelay`
- [Image list](https://gitlab.com/jonas.jasas/httprelay/container_registry)
- Run: `docker run -p 8080:8080 registry.gitlab.com/jonas.jasas/httprelay`

### Build
Install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:

```bash
go get gitlab.com/jonas.jasas/httprelay
cd ~/go/src/gitlab.com/jonas.jasas/httprelay
go run ./cmd/...
```

Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

### Test installation

Go to http://localhost:8080/health should display version number. 

## Usage examples

### Sync example (p2p)

**[Example](https://jsfiddle.net/jasajona/q6uhLuqf/)**

### Link example (p2p)
Link communication method provides peer to peer synchronous data transfers.
Data transfer is one directional, implements producer -> consumer pattern.  
Link communication method must be used when there is only one receiver and sender must know when receiver received data.

- Send data: `POST https://demo.httprelay.io/link/your_secret_channel_id`
- Receive data `GET https://demo.httprelay.io/link/your_secret_channel_id`

Sender's request will be finished when receiver makes the request.
If receiver makes request prior sender, receiver request will wait till sender makes the request.

**[Example](https://jsfiddle.net/jasajona/y35rLnd9/)**

### Multicast example (one to many)
Multicast communication method provides one to many data transfers.
Multicast communication method must be used when there are many receivers and sender don't need to know when or if receivers received it's data.

- Send data: `POST https://httprelay.io/mcast/your_secret_channel_id`
- Receive data `GET https://httprelay.io/mcast/your_secret_channel_id`

Sender's request will finish as soon as all data is transferred to the server.
If receiver makes request prior sender, receiver request will idle till sender makes the request.
Each request receiver receives cookie "SeqId" with the sequence number.
On next request receiver will wait till sender sends new data.
Cookies must be enabled on receiver side or it will receive same data multiple times and there will be no way to tell when new data is available.

- **[Message transfer example](https://jsfiddle.net/jasajona/ntwmheaf/)**
- **[Multiuser painting example](https://jsfiddle.net/jasajona/ky0cLgf9/)**



## Command line arguments
- **-a** Bind to IP address (default 0.0.0.0)
- **-p** Bind to port (default 8080)
- **-u** Bind to Unix socket path
- **-h** Print help
