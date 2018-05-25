package lightning

import (
  "fmt"
  "log"
  "path/filepath"
  "io/ioutil"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  viper "github.com/spf13/viper"
  macaroon "gopkg.in/macaroon.v2"
  lnrpc "github.com/lightningnetwork/lnd/lnrpc"
  macaroons "github.com/lightningnetwork/lnd/macaroons"
)


var (
  defaultTLSCertFilename string
  defaultMacaroonFilename string
  defaultLndDir string
  defaultRPCServer string
	defaultTLSCertPath string
	defaultMacaroonPath string
)

func GetInfo() (e error) {
  ctxb := context.Background()
  client, cleanUp := getClient()
  defer cleanUp()

  req := &lnrpc.GetInfoRequest{}
  resp, err := client.GetInfo(ctxb, req)

  if err != nil {
    log.Fatalf("failed to get info: %v", err)
  }
  fmt.Printf("%v", resp)
  return nil
}

func GetInvoice(amt int64, memo string) (string, error) {
	client, cleanUp := getClient()
	defer cleanUp()

  invoice := &lnrpc.Invoice{
		Memo:            memo,
		Value:           amt,
	}

	resp, err := client.AddInvoice(context.Background(), invoice)
	if err != nil {
			panic(err)
	}

	return resp.PaymentRequest, nil
}

func loadConfig() {
  viper.SetConfigName("config") // name of config file (without extension)
  viper.AddConfigPath("$HOME/.navybluesilver")
  err := viper.ReadInConfig()
  if err != nil {
	   panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }
  defaultTLSCertFilename  = viper.GetString("lightning.defaultTLSCertFilename")
  defaultMacaroonFilename = viper.GetString("lightning.defaultMacaroonFilename")
  defaultLndDir 					= viper.GetString("lightning.defaultLndDir")
  defaultRPCServer        = viper.GetString("lightning.defaultRPCServer")
  defaultTLSCertPath  = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath = filepath.Join(defaultLndDir, defaultMacaroonFilename)
}


func getClient() (lnrpc.LightningClient, func()) {
	conn := getClientConn(false)
	cleanUp := func() {
		conn.Close()
	}
	return lnrpc.NewLightningClient(conn), cleanUp
}

func getClientConn(skipMacaroons bool) *grpc.ClientConn {
  loadConfig()

	// Load the specified TLS certificate and build transport credentials
	// with it.
	tlsCertPath := defaultTLSCertPath
	creds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		log.Fatalf("fail to dial: %v :%s", err, tlsCertPath)
	}

	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Only process macaroon credentials if --no-macaroons isn't set and
	// if we're not skipping macaroon processing.
	if !skipMacaroons {
		// Load the specified macaroon file.
		macPath := defaultMacaroonPath
		macBytes, err := ioutil.ReadFile(macPath)
		if err != nil {
			 log.Fatalf("fail to dial: %v", err)
		}
		mac := &macaroon.Macaroon{}
		if err = mac.UnmarshalBinary(macBytes); err != nil {
 			log.Fatalf("fail to dial: %v", err)
		}

		macConstraints := []macaroons.Constraint{
			// We add a time-based constraint to prevent replay of the
			// macaroon. It's good for 60 seconds by default to make up for
			// any discrepancy between client and server clocks, but leaking
			// the macaroon before it becomes invalid makes it possible for
			// an attacker to reuse the macaroon. In addition, the validity
			// time of the macaroon is extended by the time the server clock
			// is behind the client clock, or shortened by the time the
			// server clock is ahead of the client clock (or invalid
			// altogether if, in the latter case, this time is more than 60
			// seconds).
			macaroons.TimeoutConstraint(60),

			// Lock macaroon down to a specific IP address.
			//macaroons.IPLockConstraint(ctx.GlobalString("macaroonip")),
		}

		// Apply constraints to the macaroon.
		constrainedMac, err := macaroons.AddConstraints(mac, macConstraints...)
		if err != nil {
 			log.Fatalf("fail to dial: %v", err)
		}

		// Now we append the macaroon credentials to the dial options.
		cred := macaroons.NewMacaroonCredential(constrainedMac)
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	}

	conn, err := grpc.Dial(defaultRPCServer, opts...)
	if err != nil {
 		log.Fatalf("fail to dial: %v", err)
	}
	return conn
}
