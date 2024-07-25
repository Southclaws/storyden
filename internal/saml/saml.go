package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/url"
	"os"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/crewjam/saml/samlsp"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

type SAML struct {
	Middleware *samlsp.Middleware
}

func newSAML(ctx context.Context, cfg config.Config) (*SAML, error) {
	if !cfg.SAMLEnabled {
		return nil, nil
	}

	certFile := cfg.SAMLCertFile
	keyFile := cfg.SAMLKeyFile
	idpMetaFile := cfg.SAMLIDPMetaFile

	keyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to open certificate and key files"))
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to parse certificate"))
	}

	b, err := os.ReadFile(idpMetaFile)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to open IDP metadata file"))
	}

	samlED, err := samlsp.ParseMetadata(b)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to parse IDP metadata"))
	}

	// We need to use the actual backend's URL here due to the way SAML works.
	rootURL, err := url.Parse(cfg.PublicApiAddress + "/api")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to parse public web address as a valid URL"))
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: samlED,
	})
	if err != nil {
		return nil, err
	}

	// Overwrite the default SAML paths because Storyden is always under `/api`.
	samlSP.ServiceProvider.MetadataURL = *rootURL.ResolveReference(&url.URL{Path: "api/saml/metadata"})
	samlSP.ServiceProvider.AcsURL = *rootURL.ResolveReference(&url.URL{Path: "api/saml/acs"})
	samlSP.ServiceProvider.SloURL = *rootURL.ResolveReference(&url.URL{Path: "api/saml/slo"})

	return &SAML{
		Middleware: samlSP,
	}, nil
}

func Build() fx.Option {
	return fx.Provide(newSAML)
}
