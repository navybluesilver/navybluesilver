# Navy Blue Silver
navybluesilver.net

This website was generated as a proof of concept to help me understand the development required to accept lightning donations.

In addition, I will be sharing some of my opinions on bitcoin for educational purposes.


## generate you TLS keys via:

* install and sync bitcoind
* install and sync lnd
* configure and save the following configuration file $HOME/.navybluesilver/config.json
```json                                                                                                               
{
  "lightning": {
  "defaultTLSCertFilename": "tls.cert",
  "defaultMacaroonFilename": "admin.macaroon",
      "defaultLndDir": "/root/.lnd",
      "defaultRPCServer": "127.0.0.1:10009"
  },
  "web": {
  "certFile":"cert.pem",
  "keyFile":"key.pem"
  }
}
```
