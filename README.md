# HTTP Relay
Client to client HTTP communication

## Features
- Communication between two or more HTTP clients (e.g. web browsers, mobile apps, IoT devices etc.)
- Peer to peer mode.
- Multicast (many receivers) mode.
- Simple. Secure. Blazing fast.

## Download
[Download executables](https://gitlab.com/jonas.jasas/httprelay/-/jobs/artifacts/master/browse/download?job=build:download)

## Usage examples

### Link example (p2p)
Link communication method provides peer to peer data transfers.
Link communication method must be used when there is only one receiver and sender must know when receiver received data.

- Send data: POST https://httprelay.io/link/your_secret_channel_id
- Receive data GET https://httprelay.io/link/your_secret_channel_id

Sender's request will be finished when receiver makes the request.
If receiver makes request prior sender, receiver request will wait till sender makes the request.

<script async src="//jsfiddle.net/jasajona/q6uhLuqf/embed/"></script>

