package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type Configuration struct {
	CACert  string
	AppCert string
	AppKey  string
}

func InitTlsConfig(configuration *Configuration) (*tls.Config, error) {
	appPrivateKeyRaw, err := os.ReadFile(configuration.AppKey)
	if err != nil {
		return nil, fmt.Errorf("filed to load private key, %v", err)
	}
	appCertificatRaw, err := os.ReadFile(configuration.AppCert)
	if err != nil {
		return nil, fmt.Errorf("filed to load app cert, %v", err)
	}
	appCertificate, err := tls.X509KeyPair(appCertificatRaw, appPrivateKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("filed to init X509 cert, %v", err)
	}

	caCertificateRaw, err := os.ReadFile(configuration.CACert)
	if err != nil {
		return nil, fmt.Errorf("filed to load app cert, %v", err)
	}
	rootCertpool := x509.NewCertPool()
	rootCertpool.AppendCertsFromPEM(caCertificateRaw)

	result := tls.Config{
		Certificates:     []tls.Certificate{appCertificate},
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP521},

		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			for _, a := range verifiedChains {
				for _, b := range a {
					fmt.Printf("Peer Cerr: %s\n", b.Subject.CommonName)
				}
			}
			return nil
		},

		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  rootCertpool,
		RootCAs:    rootCertpool,
	}

	return &result, nil
}
