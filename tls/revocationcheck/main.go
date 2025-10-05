package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/crypto/ocsp"
)

const (
	RevocationStatusGood RevocationStatus = iota
	RevocationStatusRevoked
	RevocationStatusUnknown
)

type RevocationStatus int

func (r RevocationStatus) String() string {
	switch r {
	case RevocationStatusGood:
		return "Good"
	case RevocationStatusRevoked:
		return "Revoked"
	case RevocationStatusUnknown:
		return "Unknown"
	default:
		return "Invalid status"
	}
}

// This simple program demonstrates how to check the revocation status of a TLS certificate
// using: CRL, OCSP, and OCSP Stampling. It connects to a specified server, retrieves the
// certificate, and checks its revocation status.
// The default Go TLS client does not perform revocation checks, so this example implements them manually.
//
// Usage: go run main.go <url>
//
// There is a site with valid, expired, and revoked certificates for testing purposes:
// https://www.ssl.com/sample-valid-revoked-and-expired-ssl-tls-certificates/
func main() {
	address := os.Args[1]

	u, err := url.Parse(address)
	if err != nil {
		log.Fatalf("invalid address %s: %v", address, err)
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
		switch u.Scheme {
		case "https":
			port = "443"
		case "ftps":
			port = "990"
		default:
			log.Fatalf("cannot infere port from address %s, specify port explicitly", address)
		}
	}

	address = net.JoinHostPort(host, port)

	tlsConfig := tls.Config{
		VerifyConnection: func(state tls.ConnectionState) error {
			if len(state.VerifiedChains) == 0 {
				return fmt.Errorf("certificate verification is disabled")
			}

			cert := state.VerifiedChains[0][0]
			issuer := state.VerifiedChains[0][1]
			var isGood bool

			fmt.Println("Checking revocation status of the certificate:", cert.SerialNumber)
			fmt.Println("Issued to:", cert.Subject)
			fmt.Println("Issued by:", cert.Issuer)
			fmt.Println()

			crlStatus, err := crlRevocationCheck(cert)
			if err != nil {
				fmt.Println("Cannot check revocation status via CRL:", err)
			} else {
				fmt.Println("CRL revocation status:", crlStatus)
				isGood = crlStatus == RevocationStatusGood
			}
			fmt.Println()

			ocspStatus, err := ocspRevocationCheck(cert, issuer)
			if err != nil {
				fmt.Println("Cannot check revocation status via OCSP:", err)
			} else {
				fmt.Println("OCSP revocation status:", ocspStatus)
				isGood = isGood || ocspStatus == RevocationStatusGood
			}
			fmt.Println()

			ocspStamplingStatus, err := parseOcspStatus(state.OCSPResponse, issuer)
			if err != nil {
				fmt.Println("Cannot check revocation status via OCSP Stampling:", err)
			} else {
				fmt.Println("Stapled OCSP revocation status:", ocspStamplingStatus)
				isGood = isGood || ocspStamplingStatus == RevocationStatusGood
			}
			fmt.Println()

			if !isGood {
				fmt.Println("Certificate has been revoked or its status is unknown")
			} else {
				fmt.Println("Certificate is valid")
			}
			return nil
		},
	}

	conn, err := tls.Dial("tcp", address, &tlsConfig)
	if err != nil {
		log.Fatalf("unable to establish TLS connection with %s: %v", address, err)
	}
	defer conn.Close()
}

func ocspRevocationCheck(cert, issuer *x509.Certificate) (RevocationStatus, error) {
	if len(cert.OCSPServer) == 0 {
		return RevocationStatusUnknown, fmt.Errorf("no OCSP server found in certificate")
	}

	ocspReq, err := ocsp.CreateRequest(cert, issuer, nil)
	if err != nil {
		return RevocationStatusUnknown, fmt.Errorf("unable to creace ocsp request: %v", err)
	}

	var ocspErrs []error
	for _, ocspServer := range cert.OCSPServer {
		fmt.Printf("Requesting status from OCSP server %s\n", ocspServer)
		resp, err := http.Post(ocspServer, "application/ocsp-request", bytes.NewReader(ocspReq))
		if err != nil {
			ocspErrs = append(ocspErrs, fmt.Errorf("error while requesting ocsp server %s: %v", ocspServer, err))
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ocspErrs = append(ocspErrs, fmt.Errorf("unable to get body from ocsp server %s response: %v", ocspServer, err))
			continue
		}
		ocspStatus, err := parseOcspStatus(body, issuer)
		if err != nil {
			ocspErrs = append(ocspErrs, fmt.Errorf("invalid ocsp response from server %s: %v", ocspServer, err))
			continue
		}
		return ocspStatus, nil
	}
	return RevocationStatusUnknown, errors.Join(ocspErrs...)
}

func parseOcspStatus(ocspRespRaw []byte, issuer *x509.Certificate) (RevocationStatus, error) {
	if len(ocspRespRaw) == 0 {
		return RevocationStatusUnknown, fmt.Errorf("no OCSP response provided")
	}
	ocspResp, err := ocsp.ParseResponse(ocspRespRaw, issuer)
	if err != nil {
		return RevocationStatusUnknown, err
	}

	switch ocspResp.Status {
	case ocsp.Good:
		return RevocationStatusGood, nil
	case ocsp.Revoked:
		return RevocationStatusRevoked, nil
	case ocsp.Unknown:
		return RevocationStatusUnknown, nil
	default:
		return RevocationStatusUnknown, fmt.Errorf("unknown ocsp status: %d", ocspResp.Status)
	}
}

func crlRevocationCheck(cert *x509.Certificate) (RevocationStatus, error) {
	if len(cert.CRLDistributionPoints) == 0 {
		return RevocationStatusUnknown, fmt.Errorf("no CRL distribution points found in certificate")
	}

	var crlErrs []error
	for _, crlURL := range cert.CRLDistributionPoints {
		fmt.Printf("Fetching CRL from %s\n", crlURL)
		resp, err := http.Get(crlURL)
		if err != nil {
			crlErrs = append(crlErrs, fmt.Errorf("error fetching CRL from %s: %v", crlURL, err))
			continue
		}
		crlRaw, err := io.ReadAll(resp.Body)
		if err != nil {
			crlErrs = append(crlErrs, fmt.Errorf("error reading CRL from %s: %v", crlURL, err))
			continue
		}

		if p, _ := pem.Decode(crlRaw); p != nil {
			crlRaw = p.Bytes
		}

		crl, err := x509.ParseRevocationList(crlRaw)
		if err != nil {
			crlErrs = append(crlErrs, fmt.Errorf("error parsing CRL from %s: %v", crlURL, err))
			continue
		}

		for _, revoked := range crl.RevokedCertificateEntries {
			if revoked.SerialNumber.Cmp(cert.SerialNumber) == 0 {
				return RevocationStatusRevoked, nil
			}
		}
		return RevocationStatusGood, nil
	}
	return RevocationStatusUnknown, errors.Join(crlErrs...)
}
