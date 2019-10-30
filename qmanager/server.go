package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"crypto/rand"
    "crypto/rsa"
    "crypto/x509"
	"encoding/pem"
	"crypto/sha512"
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

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		 
	}
	return plaintext
} 

type Args struct {
	// enc_message []byte
	message []byte
 }

type CreateExtraData struct{}

func (t *CreateExtraData) Create(args *Args, reply *string) error {
	// *reply = args.X + args.Y
	// return nil
	// priv, pub := GenerateRsaKeyPair()
	// _ = pub
	// // fmt.Println(pub)

	// decrypted_message := DecryptWithPrivateKey (args.enc_message, priv)
	fmt.Println(args)
	// fmt.Println(decrypted_message)

    // // Export the keys to pem string
    // priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
    // pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

    // // Import the keys from pem string
    // priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
    // pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

    // // Export the newly imported keys
    // priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	// pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)
	// if priv_pem != priv_parsed_pem || pub_pem != pub_parsed_pem {
	// 	fmt.Println("Failure: Export and Import did not result in same Keys")
	// } else {
	// 		fmt.Println("Success")
	// }


    // fmt.Println(priv)
    // fmt.Println(pub_pem)
	*reply = "SUCCESS"
	return nil
}


func main() {

	

    // Check that the exported/imported keys match the original keys
    
	cal := new(CreateExtraData)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}