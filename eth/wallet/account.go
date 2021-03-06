package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	version        = byte(0x0000)
	checksumLength = 4
)

type Account struct {
	PublicKey  []byte
	PrivateKey ecdsa.PrivateKey
}

func (w Account) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionedHash)
	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	fmt.Printf("pubKey is :%x\n", w.PublicKey)
	fmt.Printf("pubHash is :%x\n", pubHash)
	fmt.Printf("checksum is :%x\n", checksum)
	fmt.Printf("address is %x\n", address)
	return address
}

func PublicKeyHash(publickey []byte) []byte {
	pubHash := sha256.Sum256(publickey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:checksumLength]
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func NewAccount() *Account {
	privateKey, publicKey := NewKeyPair()
	return &Account{publicKey, privateKey}
}
