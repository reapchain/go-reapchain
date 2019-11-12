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
    "crypto"
//     "math/big"
// "encoding/base64"
    "errors"
    "os"
    "io/ioutil"
    "github.com/ethereum/go-ethereum/common"
 
)

 
type Validatorsset struct {
        Senator               defaultValidator
        Parliamentarian   defaultValidator
        Candidate         defaultValidator
        General            defaultValidator
        Qmanager         defaultValidator
}

       
type defaultValidator struct {   
        address common.Address   
        tag uint64
}

type NodeData struct {
        IPAddress       net.IP
        PortNo  uint16
        PublicKey       common.Address 
        Tag     uint64
        QRND    uint64
        
}


var secret = []byte("TbqZ6fkx2T") 
 
func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
    privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
    return privkey, &privkey.PublicKey
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

func VerifyWithPublicKey(  sig []byte, pub *rsa.PublicKey) error {
        // hash := sha512.New()
        digest := sha512.Sum512(secret)
        return rsa.VerifyPKCS1v15( pub, crypto.SHA512, digest[:], sig)
	 
	 
}

type RequestEnroll struct
{
        NodeData NodeData
        Signature []byte
}

type ExtraDataResponse struct
{
        ValidatorSet Validatorsset
        Signature []byte
}
type Args struct {
	// enc_message []byte
        Message []byte
        Signature []byte
 }

func main() {
	// fmt.Println(message)

        public_file_reader, err := ioutil.ReadFile("public_key.pem")
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        pub, err := ParseRsaPublicKeyFromPemStr(string(public_file_reader))

        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
     
	
	signature := EncryptWithPublicKey(secret, pub)
	// fmt.Println(encoded_message)
	 

	client, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
        }
        args := &Args{secret , signature}

	var reply ExtraDataResponse
	c := jsonrpc.NewClient(client)
	err = c.Call("ExtraData.Create", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
         
        Signature := &reply.Signature
        Vals := &reply.ValidatorSet
        fmt.Println(Signature)
        fmt.Println(VerifyWithPublicKey( *Signature, pub))
        VerificationCheck := VerifyWithPublicKey(*Signature, pub)
        if  VerificationCheck == nil {
                
                fmt.Println(Vals)
        }
}