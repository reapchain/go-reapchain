package quantum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"os"
	"syscall"
	"time"
)


const bufferSize = 16384 *3

func GenerateQrnd() []byte  {
	var ioldIndex byte
	var Qrnddata []byte
	for {
		buffer := readUpFile()   //buffer data type is slice []byte
		if len(buffer) == 0 {
			writeDnFile(0)
			time.Sleep(1 * time.Second)
		} else if buffer[0] == ioldIndex {
			writeDnFile(buffer[0])
			time.Sleep(1 * time.Second)
		} else {
			ioldIndex = buffer[0]
			if len(buffer) > 0 {
				writeDnFile(buffer[0])
				Qrnddata = buffer
			} else {
				time.Sleep(1 * time.Second)
			}
			break
		}
	}
	return Qrnddata
}

func readUpFile() ([]byte, error)  {
	buffer := make([]byte, bufferSize)

	flags := os.O_RDONLY | os.O_EXCL  | syscall.O_DIRECT
	upFile, err := os.OpenFile(podc_global.QRNGFilePrefix+"up.ini", flags, 0644)
	if err != nil {
		log.Error("QRND Device", "Open Error", err.Error())
		return nil
	}
	_, err = upFile.ReadAt(buffer, 0)
	if err != nil {
		log.Error("QRND Device", "Read Error", err.Error())
		return nil
	}
	upFile.Close()

	return buffer
}

func writeDnFile(data byte) {
 	os.Remove(podc_global.QRNGFilePrefix + "dn.ini")
	dnFile, err := os.OpenFile(podc_global.QRNGFilePrefix+"dn.ini", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("CREATE-ERROR: dn.ini")
		panic(err.Error())
	}
	writtenSize, err := dnFile.WriteAt([]byte{data}, 0)
	if err != nil || writtenSize == 0 {
		fmt.Println("WRITE-ERROR: dn.ini")
		panic(err.Error())
	}
	dnFile.Close()
}

