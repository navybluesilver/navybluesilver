# navybluesilver.net

This website was generated as a proof of concept to help me understand the development required to accept lightning donations.

In addition, I will be sharing some of my opinions on bitcoin for educational purposes.

## Installation

### Preliminaries
#### install and sync bitcoind
[`bitcoind`](https://github.com/bitcoin/bitcoin)

#### install and sync lnd
[`lnd`](https://github.com/lightningnetwork/lnd)

#### download and build navybluesilver
```
$ cd $GOPATH/src/github.com/
$ git clone https://github.com/navybluesilver/navybluesilver.git
$ cd $GOPATH/src/github.com/navybluesilver
$ go build
```
#### configure
```
$ nano /home/navybluesilver/.navybluesilver/config.json
```

```json                                                                                                               
{
  "lightning": {
  "defaultTLSCertFilename": "tls.cert",
  "defaultMacaroonFilename": "admin.macaroon",
  "defaultLndDir": "/home/navybluesilver/.navybluesilver",
  "defaultRPCServer": "127.0.0.1:10009"
  },
  "web": {
  "certFile":"/home/navybluesilver/.navybluesilver/fullchain.pem",
  "keyFile":"/home/navybluesilver/.navybluesilver/privkey.pem"
  }

}
```
#### certificates
Generate TLS certificates (e.g. by using [letsencrypt](https://letsencrypt.org/))
Copy TLS and .lnd certificates to /home/navybluesilver/.navybluesilver/
