package quantum

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const bufferSize = 16384

const filePrefix ="/Volumes/PSoC USB/"  //on my MAC

func GenerateQrnd() []byte  {
	var ioldIndex byte
	var Qrnddata []byte


	for {
		buffer := readUpFile()   //buffer data type is slice []byte
		if buffer[0] == ioldIndex {
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

func readUpFile() []byte {
	buffer, err := ioutil.ReadFile(filePrefix + "up.ini")
	if err != nil {
		fmt.Println("OPEN-ERROR: up.ini")
		panic(err.Error())
	}

	return buffer
}

func writeDnFile(data byte) {
	os.Remove(filePrefix + "dn.ini")
	err := ioutil.WriteFile(filePrefix+"dn.ini", []byte{data}, 0644)
	if err != nil {
		fmt.Println("CREATE-ERROR: dn.ini")
		panic(err.Error())
	}
}

