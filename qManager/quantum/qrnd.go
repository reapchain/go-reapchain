package quantum

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"os"
	"syscall"
	"time"
)


const bufferSize = 16384 *3

func GenerateQRNGData() ([]byte, error)  {
	var ioldIndex byte
	var Qrngdata []byte
	for {
		buffer, err := readUpFile()   //buffer data type is slice []byte
		if err == nil{

			if len(buffer) == 0 {
				writeErr := writeDnFile(0)
				if writeErr == nil {
					time.Sleep(1 * time.Second)
				}else{
					return nil, writeErr
				}
			} else if buffer[0] == ioldIndex {
				writeErr := writeDnFile(buffer[0])
				if writeErr == nil {
					time.Sleep(1 * time.Second)
				}else{
					return nil, writeErr
				}
			} else {
				ioldIndex = buffer[0]
				if len(buffer) > 0 {
					writeErr := writeDnFile(buffer[0])
					if writeErr == nil {
						Qrngdata = buffer
					}else{
						return nil, writeErr
					}
				} else {
					time.Sleep(1 * time.Second)
				}
				break
			}
		} else{
			return nil, err
		}

	}
	return Qrngdata, nil
}

func readUpFile() ([]byte, error)  {
	buffer := make([]byte, bufferSize)

	flags := os.O_RDONLY | os.O_EXCL  | syscall.O_DIRECT
	upFile, err := os.OpenFile(podc_global.QRNGFilePrefix+"up.ini", flags, 0644)
	if err != nil {
		log.Error("QRNG Device", "Open Error", err.Error())
		return nil, err
	}
	_, err = upFile.ReadAt(buffer, 0)
	if err != nil {
		log.Error("QRNG Device", "Read Error", err.Error())
		return nil, err
	}
	upFile.Close()

	return buffer, nil
}

func writeDnFile(data byte) error {
 	os.Remove(podc_global.QRNGFilePrefix + "dn.ini")
	dnFile, err := os.OpenFile(podc_global.QRNGFilePrefix+"dn.ini", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("QRNG Device", "Create Error", err.Error())
		return err
	}

	writtenSize, err := dnFile.WriteAt([]byte{data}, 0)
	if err != nil || writtenSize == 0 {
		log.Error("QRNG Device", "Write Error", err.Error())
		return err
	}

	dnFile.Close()
	return nil
}

