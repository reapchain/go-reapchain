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
        "crypto"
        "crypto/sha512"
        // "encoding/base64"
        // "crypto/sha1"
    "errors"
    "os"
    "io/ioutil"
    "github.com/ethereum/go-ethereum/common"
    "github.com/patrickmn/go-cache"
    "time"


)
var go_cache = cache.New(5*time.Minute, 10*time.Minute)
var secret = []byte("TbqZ6fkx2T") 

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

func InitializeKeys()(bool){
        priv, pub := GenerateRsaKeyPair()
        fmt.Println(priv)
        fmt.Println(pub)
        pemPrivateFile, err := os.Create("private_key.pem")
                

        var pemPrivateBlock = &pem.Block{
                Type:  "RSA PRIVATE KEY",
                Bytes: x509.MarshalPKCS1PrivateKey(priv),
                }

                err = pem.Encode(pemPrivateFile, pemPrivateBlock)
                if err != nil {
                // fmt.Println(err)
                // os.Exit(1)
                return false
                }

                pemPrivateFile.Close()


                pemPublicFile, err := os.Create("public_key.pem")

                pubkey_bytes, err := x509.MarshalPKIXPublicKey(pub)

                var pemPublickBlock = &pem.Block{
                        Type:  "RSA PUBLIC KEY",
                        Bytes: pubkey_bytes,
                }

                err = pem.Encode(pemPublicFile, pemPublickBlock)
                if err != nil {
                        // fmt.Println(err)
                        // os.Exit(1)
                        return false
                }

                pemPublicFile.Close()
                return true

}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		 
	}
	return ciphertext
}



// func EncryptWithPrivateKey(msg []byte, priv *rsa.PrivateKey) []byte {
// 	hash := sha512.New()
// 	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, priv, msg, nil)
// 	if err != nil {
		 
// 	}
// 	return ciphertext
// }

func SignWithPrivateKey(msg []byte, priv *rsa.PrivateKey) (signature []byte, err error) {
        // hash := sha512.New()
        digest := sha512.Sum512(msg)
        signature, signErr := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA512, digest[:])
	 
	return signature, signErr
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) (signature []byte, err error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	 
	return plaintext, err
} 

type Args struct {
	// enc_message []byte
        Message []byte
        Signature []byte
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

type ParseMessage struct{}
type ExtraData struct{}
type Enrollment struct{}

func (t *Enrollment) Enroll(req_args *RequestEnroll, reply *Args) error {
	// *reply = args.X + args.Y
        // return nil
        
        private_file_reader, err := ioutil.ReadFile("private_key.pem")
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        priv, err := ParseRsaPrivateKeyFromPemStr(string(private_file_reader))

        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

       

        decrypted_message, err := DecryptWithPrivateKey (req_args.Signature, priv)
        
        if err != nil{
                fmt.Println(err)
                os.Exit(1)
        }
	// fmt.Println(message)
        // fmt.Println(decrypted_message)
        decrypted_string := string(decrypted_message)
        fmt.Println(decrypted_string)

	 
        go_cache.Set("TMP_NODE", req_args.NodeData, cache.DefaultExpiration)
        foo, found := go_cache.Get("TMP_NODE")
	if found {
		fmt.Println(foo)
	}

        
        // message_test := []byte("private_test")
        // encoded_message := EncryptWithPrivateKey(message_test, priv)
        // fmt.Println(encoded_message)
        
         
 
        signature, signErr := SignWithPrivateKey(secret, priv)

        if signErr != nil {
		fmt.Println("Could not sign message:%s", signErr.Error())
        }
        
       
        // fmt.Println(signature)
        // fmt.Println(args)
	*reply = *&Args{secret, signature}
	return nil
}

func (t *ExtraData) Create(args *Args, reply *ExtraDataResponse) error {
	// *reply = args.X + args.Y
        // return nil
        
        private_file_reader, err := ioutil.ReadFile("private_key.pem")
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        priv, err := ParseRsaPrivateKeyFromPemStr(string(private_file_reader))

        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

       

        decrypted_message, err := DecryptWithPrivateKey (args.Signature, priv)
        
        if err != nil{
                fmt.Println(err)
                os.Exit(1)
        }
	// fmt.Println(message)
        // fmt.Println(decrypted_message)
        decrypted_string := string(decrypted_message)
        fmt.Println(decrypted_string)

        sen_temp, _ := go_cache.Get("Senator")
        Senator := sen_temp.(defaultValidator)

        par_temp, _ := go_cache.Get("Parliamentarian")
        Parliamentarian, _ := par_temp.(defaultValidator)

        can_temp, _ := go_cache.Get("Candidate")
        Candidate, _ := can_temp.(defaultValidator)

        gen_temp, _ := go_cache.Get("General")
        General, _ := gen_temp.(defaultValidator)

        qm_temp, _ := go_cache.Get("Qmanager")
        Qmanager, _ := qm_temp.(defaultValidator)

        // fmt.Println("SENATOR")

        // fmt.Println(Senator)
        temp_val_set := Validatorsset{Senator: Senator, Parliamentarian: Parliamentarian, Candidate: Candidate, General: General, Qmanager: Qmanager}

        signature, signErr := SignWithPrivateKey(secret, priv)

        if signErr != nil {
		fmt.Println("Could not sign message:%s", signErr.Error())
        }
        
        
         
        // go_cache.Set("Senator", temp_senator, cache.DefaultExpiration)
        // go_cache.Set("Parliamentarian", temp_parliamentary, cache.DefaultExpiration)
        // go_cache.Set("Candidate", temp_candidate, cache.DefaultExpiration)
        // go_cache.Set("General", temp_general, cache.DefaultExpiration)
        // go_cache.Set("Qmanager", temp_qmanager, cache.DefaultExpiration)

        fmt.Println(ExtraDataResponse{ValidatorSet: temp_val_set, Signature: signature})

	*reply = *&ExtraDataResponse{ValidatorSet: temp_val_set, Signature: signature}
	return nil
}

func (t *ParseMessage) Create(args *Args, reply *Args) error {
	// *reply = args.X + args.Y
        // return nil
        
        private_file_reader, err := ioutil.ReadFile("private_key.pem")
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        priv, err := ParseRsaPrivateKeyFromPemStr(string(private_file_reader))

        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        // public_file_reader, err := ioutil.ReadFile("public_key.pem")
        // if err != nil {
        //         fmt.Println(err)
        //         os.Exit(1)
        // }
        // pub, err := ParseRsaPublicKeyFromPemStr(string(public_file_reader))

        // if err != nil {
        //         fmt.Println(err)
        //         os.Exit(1)
        // }
 

	decrypted_message, err := DecryptWithPrivateKey (args.Signature, priv)
        
        if err != nil{
                fmt.Println(err)
                os.Exit(1)
        }
	// fmt.Println(message)
        fmt.Println(decrypted_message)
        decrypted_string := string(decrypted_message)
        fmt.Println(decrypted_string)
        
        // message_test := []byte("private_test")
        // encoded_message := EncryptWithPrivateKey(message_test, priv)
        // fmt.Println(encoded_message)
        
        data := []byte("Private_TEST")
 
        signature, signErr := SignWithPrivateKey(data, priv)

        if signErr != nil {
		fmt.Println("Could not sign message:%s", signErr.Error())
        }
        
        // fmt.Println((signature))

        // b64sig := base64.StdEncoding.EncodeToString(signature)

	// decodedSignature, _ := base64.StdEncoding.DecodeString(b64sig)

	// // verify part

	// verifyErr := rsa.VerifyPKCS1v15(pub, crypto.SHA512, digest[:], decodedSignature)

	// if verifyErr != nil {
	// 	fmt.Println("Verification failed: %s", verifyErr)
        // }
        // fmt.Println(verifyErr)

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
        args = &Args{data, signature}
        fmt.Println(data)
        fmt.Println(signature)
        fmt.Println(args)
	*reply = *args
	return nil
}



func main() {

	if _, err := os.Stat("private_key.pem"); os.IsNotExist(err) {
                InitializeKeys()

        }    

        temp_senator := defaultValidator{address: common.HexToAddress("0xdceceaf3fc5c0a63d195d69b1a90011b7b19650d"), tag: 11154013587666973726}
        temp_parliamentary := defaultValidator{address: common.HexToAddress("0x7d577a597b2742b498cb5cf0c26cdcd726d39e6e"), tag: 17349330445822630798}
        temp_candidate := defaultValidator{address: common.HexToAddress("0x598443f1880ef585b21f1d7585bd0577402861e5"), tag: 5577006791947779410}
        temp_general := defaultValidator{address: common.HexToAddress("0x13cbb8d99c6c4e0f2728c7d72606e78a29c4e224"), tag: 2117424165630711375}
        temp_qmanager := defaultValidator{address: common.HexToAddress("0x77db2bebba79db42a978f896968f4afce746ea1f"), tag: 11123412387666973726}
        fmt.Println(temp_senator)

        go_cache.Set("Senator", temp_senator, cache.DefaultExpiration)
        go_cache.Set("Parliamentarian", temp_parliamentary, cache.DefaultExpiration)
        go_cache.Set("Candidate", temp_candidate, cache.DefaultExpiration)
        go_cache.Set("General", temp_general, cache.DefaultExpiration)
        go_cache.Set("Qmanager", temp_qmanager, cache.DefaultExpiration)



        cal := new(ParseMessage)
        cal2 := new(ExtraData)
        enroll := new(Enrollment)
	server := rpc.NewServer()
        server.Register(cal)
        server.Register(cal2)
        server.Register(enroll)
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