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

### Multicast example (one to many)
Multicast communication method provides one to many data transfers.
Multicast communication method must be used when there are many receivers and sender don't need to know when or if receivers received it's data.

- Send data: POST https://httprelay.io/mcast/your_secret_channel_id
- Receive data GET https://httprelay.io/mcast/your_secret_channel_id

Sender's request will finish as soon as all data is transferred to the server.
If receiver makes request prior sender, receiver request will idle till sender makes the request.
Each request receiver receives cookie "SeqId" with the sequence number.
On next request receiver will wait till sender sends new data.
Cookies must be enabled on receiver side or it will receive same data multiple times and there will be no way to tell when new data is available.