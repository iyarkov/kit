package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

var ErrorUnknownPeer = errors.New("unknown peer d")

type Configuration struct {
	CACert     string
	AppCert    string
	AppKey     string
	KnownPeers []string
}

func (configuration *Configuration) NewCryptoTlsConfig() (*tls.Config, error) {
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

	peersSet := make(map[string]bool, len(configuration.KnownPeers))
	for _, peer := range configuration.KnownPeers {
		peersSet[peer] = true
	}

	result := tls.Config{
		Certificates:     []tls.Certificate{appCertificate},
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP521},
		ClientAuth:       tls.RequireAndVerifyClientCert,
		ClientCAs:        rootCertpool,
		RootCAs:          rootCertpool,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			var peerName string
			if len(verifiedChains) > 0 && len(verifiedChains[0]) > 0 {
				peerName = verifiedChains[0][0].Subject.CommonName
			}
			if peersSet[peerName] {
				zerolog.Ctx(context.Background()).Debug().Msgf("connection from known peer %s", peerName)
				return nil
			} else {
				zerolog.Ctx(context.Background()).Warn().Msgf("attempt to connect from unknown peer %s", peerName)
				return ErrorUnknownPeer
			}
		},
	}

	return &result, nil
}

type ConnectionInfo struct {
	Peer          string
	CertificateId uint64
}

func (ci *ConnectionInfo) FromPeerCertificate(peerCert *x509.Certificate) {
	ci.Peer = peerCert.Subject.CommonName
	ci.CertificateId = peerCert.SerialNumber.Uint64()
}
