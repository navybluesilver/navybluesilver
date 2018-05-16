package lightning

import (
  "fmt"
  "log"
  "strings"
  "path/filepath"
  "io/ioutil"
	"os"
	"os/user"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  macaroon "gopkg.in/macaroon.v2"
  lnrpc "github.com/lightningnetwork/lnd/lnrpc"
  macaroons "github.com/lightningnetwork/lnd/macaroons"
)

const (
	defaultTLSCertFilename  = "tls.cert"
	defaultMacaroonFilename = "admin.macaroon"
	defaultLndDir 					= "/root/.lnd"
)

var (
	defaultTLSCertPath  = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath = filepath.Join(defaultLndDir, defaultMacaroonFilename)
	defaultRPCServer    = "localhost:10009"
	webApplicationPort  = ":8081"
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

func getClient() (lnrpc.LightningClient, func()) {
	conn := getClientConn(false)
	cleanUp := func() {
		conn.Close()
	}
	return lnrpc.NewLightningClient(conn), cleanUp
}

func getClientConn(skipMacaroons bool) *grpc.ClientConn {

	// Load the specified TLS certificate and build transport credentials
	// with it.
	tlsCertPath := defaultTLSCertPath
	creds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Only process macaroon credentials if --no-macaroons isn't set and
	// if we're not skipping macaroon processing.
	if !skipMacaroons {
		// Load the specified macaroon file.
		macPath := cleanAndExpandPath(defaultMacaroonPath)
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

func cleanAndExpandPath(path string) string {
	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {
		var homeDir string

		user, err := user.Current()
		if err == nil {
			homeDir = user.HomeDir
		} else {
			homeDir = os.Getenv("HOME")
		}

		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%,
	// but the variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}
