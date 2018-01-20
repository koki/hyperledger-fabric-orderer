package orderer

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	signx "github.com/hyperledger/fabric/bccsp/signer"
	cryptox "github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"

	"github.com/koki/hyperledger-fabric-orderer/config"
)

type signer struct {
	cert         *x509.Certificate
	pk           bccsp.Key
	signer       crypto.Signer
	hashFunction config.IdentityIdentifierHashFunction
	mspId        string
}

func NewSignerFromConfig(config *config.OrdererConfig) (*signer, error) {
	if config == nil {
		return nil, fmt.Errorf("Obtained nil config")
	}

	if len(config.KeyStore) == 0 || len(config.SignCert) == 0 {
		return nil, fmt.Errorf("Empty Public and/or Private MSP info")
	}

	pemCert, _ := pem.Decode(config.SignCert)
	if pemCert == nil {
		return nil, fmt.Errorf("Could not decode sign-cert")
	}

	cert, err := x509.ParseCertificate(pemCert.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse sign-cert %v", err)
	}

	mbccsp := factory.GetDefault()

	certPubKey, err := mbccsp.KeyImport(cert, &bccsp.X509PublicKeyImportOpts{Temporary: true})
	if err != nil {
		return nil, fmt.Errorf("Failed to import sign-cert to correct format %v", err)
	}

	hashOpt, err := bccsp.GetHashOpt(string(config.IdentityIdentifierHashFunction))
	if err != nil {
		return nil, fmt.Errorf("Failed getting hash function options %v", err)
	}

	_, err = mbccsp.Hash(cert.Raw, hashOpt)
	if err != nil {
		return nil, fmt.Errorf("Failed hashing raw certificate to compute the id of the IdentityIdentifier %v", err)
	}

	pemPrivateKey, _ := pem.Decode(config.KeyStore)
	if pemPrivateKey == nil {
		return nil, fmt.Errorf("Could not decode key-store %v", err)
	}

	certPrivKey, err := mbccsp.KeyImport(pemPrivateKey.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: true})
	if err != nil {
		return nil, fmt.Errorf("Could not import key-store to correct format %v", err)
	}

	peerSigner, err := signx.New(mbccsp, certPrivKey)
	if err != nil {
		return nil, fmt.Errorf("Failed initializing key-store signer %v", err)
	}

	signr := &signer{
		cert:         cert,
		pk:           certPubKey,
		signer:       peerSigner,
		hashFunction: config.IdentityIdentifierHashFunction,
		mspId:        config.MspId,
	}

	return signr, nil
}

func (signr *signer) NewSignatureHeader() (*common.SignatureHeader, error) {
	pb := &pem.Block{Bytes: signr.cert.Raw}
	pemBytes := pem.EncodeToMemory(pb)
	if pemBytes == nil {
		return nil, fmt.Errorf("Encoding of identitiy failed")
	}

	sId := &msp.SerializedIdentity{Mspid: signr.mspId, IdBytes: pemBytes}
	idBytes, err := proto.Marshal(sId)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal a SerializedIdentity structure for identity %s, err %s", signr.mspId, err)
	}

	nonce, err := cryptox.GetRandomNonce()
	if err != nil {
		return nil, fmt.Errorf("Failed creating nonce [%s]", err)
	}

	sh := &common.SignatureHeader{}
	sh.Creator = idBytes
	sh.Nonce = nonce

	return sh, nil
}

func (signr *signer) Sign(message []byte) ([]byte, error) {
	hashOpt, err := bccsp.GetHashOpt(bccsp.SHA256)
	if err != nil {
		return nil, err
	}

	if signr.hashFunction == config.SHA3_256 || signr.hashFunction == config.SHA3_384 {
		hashOpt, err = bccsp.GetHashOpt(bccsp.SHA3_256)
		if err != nil {
			return nil, err
		}
	}

	mbccsp := factory.GetDefault()

	digest, err := mbccsp.Hash(message, hashOpt)
	if err != nil {
		return nil, fmt.Errorf("Failed computing digest [%s]", err)
	}

	// Sign
	return signr.signer.Sign(rand.Reader, digest, nil)
}
