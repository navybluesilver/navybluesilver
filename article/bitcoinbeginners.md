# Bitcoin for beginners

## Rule Number 1: Control your private keys

### Where do we store bitcoins?
First you have to create a public key / private key. The public key resembles an account number that you can share with other people. The private key is the password to move all the funds located on your public key. This should only be known by you. It is possible to take backups of your private key, but this increases the risk that somebody can steal your private key(s). If the private key is stored on an exchange, then your bitcoin can be censored, confiscated or hacked/stolen. It is not really *your* bitcoin. One of the main design principles behind bitcoin is to remove custodial risk (by controlling your own private key).

*If you own the private key it is your bitcoin, if you do not own the private key it is not your bitcoin*

#### Option 1: Hardware Wallet
The best way to store bitcoins is in a hardware wallet. This is also known as "cold storage". This device cost usually around 100 USD. Lately they have been sold out, so expect to wait more than a month.
[trezor.io](https://trezor.io/)

#### Option 2: Paper Wallet
A cheaper option is to use a paper wallet, [bitaddress.org](https://www.bitaddress.org/). Note that your device (windows/apple) might be infected by mall-ware, which could compromise your private key generation. In addition a hardware wallet is a lot more user friendly.

#### Option 3: Mobile App
Your mobile device will be more secure than your Apple/Windows PC. Use an app such as [samurai wallet](https://samouraiwallet.com/) to generate public keys and store private keys.

### Where do we buy bitcoins?
After creating your own public key, you want to send bitcoins to this address.

#### Option 1: Localbitcoins with cash
If you want perfect anonymity you could arrange a meeting in person with somebody in a (safe) public place. You can use [localbitcoins.net](https://localbitcoins.net/) to exchange your cash for his/her bitcoin. The website [localbitcoins.net](https://localbitcoins.net/) can also be used to buy bitcoin using a bank transfer, etc. However, this would require you to share your identity.

#### Option 2: ATM
An bitcoin ATM allows you to change cash into bitcoins. Usually the rates are more expensive than the other options.
[coinatmradar.com](https://coinatmradar.com/)

#### Option 3: decentralized exchanges (e.g. BISQ)
Bisq is an open-source desktop application that allows you to buy and sell bitcoins in exchange for national currencies, or alternative crypto currencies. Unlike traditional online exchanges, Bisq is designed to be:

* Instantly accessible – no need for registration or approval from a central authority.
* Decentralized – there is no single point of failure. The system is peer-to-peer and trading cannot be stopped or censored.
* Safe – Bisq never holds your funds. Decentralized arbitration system and security deposits protect traders.
* Private – no one except trading partners exchange personally identifying data. All personal data is stored locally.
* Secure – end-to-end encrypted communication routed over Tor.
* Open – every aspect of the project is transparent. The code is open source.
* Easy – we take usability seriously.

The advantage is that your transaction is not stored in a centralized database (managed by the exchange). The disadvantage is that you would need at least some bitcoin to do your first transaction (which is used as a deposit).

[bisq](https://bisq.network/)

#### Option 4: centralized exchange
If you don't mind sharing your identification with a third party, you can also use an exchange.

## Rule Number 2: Don't trust, verify
### Run a full node
The best way to verify your transaction is by running a [bitcoin full node](https://bitcoin.org/en/download) yourself. However, you could also "trust" a block explorer, such as [oxt.me](https://oxt.me/). Note that we verify all bitcoin transactions, not just our own. Because we are able to verify the entire global ledger ourselves, we do not need to trust any third party to ensure that there is no manual manipulation or double spends happening. If you have a technical background, you should run a full node.

*Don't trust, verify*


### Rules without rulers
To say that bitcoin is an unregulated market is false. It is true that the protocol cannot be regulated by authority. This means that bitcoin cannot be censored or confiscated. In addition the monetary policy is hard-coded in the protocol. These cannot be changed by force. However, in theory anyone can run any bitcoin software. If you change the code behind your full node, you will fork off the network. Meaning you have created an altcoin that shares a history with the original bitcoin chain. The value of your altcoin is determined by the market. Most likely that value will be 0, unless you can convince somebody to buy it. In principle bitcoin is whatever you want it to be. It is a system of *rules without rules*.

### Nation state regulation
Please read the disclaimer, this is not legal advice. The regulation around bitcoin is still vague and unclear.
The bitcoin SPOT market is not regulated by authority. This means that if the price of bitcoin is 0, there will be no bailout or refund. There is also no "cancel transaction" option. There is a regulated futures market for bitcoin:

* [CME Futures](http://www.cmegroup.com/trading/bitcoin-futures.html)
* [CBOE Futures](http://www.cboe.com/delayedquote/futures-quotes)

The current understanding is that bitcoin is more like a commodity, not a security. Note that ICO's are considered as securities:

* [SEC press release](https://www.sec.gov/news/press-release/2017-131)
* [PBC press release](http://www.pbc.gov.cn/english/130721/3377816/index.html)
