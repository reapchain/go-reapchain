// client.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"crypto/rand"
    "crypto/rsa"
	"crypto/x509"
	"crypto/sha512"
	_ "reflect"
    "encoding/pem"
    "errors"
	 
)
func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
    privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
    return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
    privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
    privkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PRIVATE KEY",
                    Bytes: privkey_bytes,
            },
    )
    return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(privPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
    pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
    if err != nil {
            return "", err
    }
    pubkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PUBLIC KEY",
                    Bytes: pubkey_bytes,
            },
    )

    return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pubPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    switch pub := pub.(type) {
    case *rsa.PublicKey:
            return pub, nil
    default:
            break // fall through
    }
    return nil, errors.New("Key type is not RSA")
}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		 
	}
	return ciphertext
}


type Args struct {
	// enc_message []byte
	message []byte

}

func main() {
	message := []byte("test")
	// fmt.Println(message)

	priv, pub := GenerateRsaKeyPair()
	_ = priv

    // Export the keys to pem string
    // priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
    // pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

    // Import the keys from pem string
    // priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
    // pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

    // // Export the newly imported keys
    // priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
    // pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

    // fmt.Println(priv_parsed_pem)
	// fmt.Println(pub_parsed_pem)
	
	encoded_message := EncryptWithPublicKey(message, pub)
	fmt.Println(encoded_message)
	// s := string(encoded_message )
	// fmt.Println(s) // ABC€

	// t := []byte(s )
	// fmt.Println(t) // ABC€

	client, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	// args := &Args{encoded_message , 1}
	args := &Args{s}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("CreateExtraData.Create", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Results: %s\n", reply)
}