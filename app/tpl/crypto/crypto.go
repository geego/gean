package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gostores/assist"
)

// New returns a new instance of the crypto-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "crypto" namespace.
type Namespace struct{}

// MD5 hashes the given input and returns its MD5 checksum.
func (ns *Namespace) MD5(in interface{}) (string, error) {
	conv, err := assist.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// SHA1 hashes the given input and returns its SHA1 checksum.
func (ns *Namespace) SHA1(in interface{}) (string, error) {
	conv, err := assist.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// SHA256 hashes the given input and returns its SHA256 checksum.
func (ns *Namespace) SHA256(in interface{}) (string, error) {
	conv, err := assist.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}
