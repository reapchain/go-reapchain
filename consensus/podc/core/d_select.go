package core

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/qManager/podc_global"

	//"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"os"

	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/log"
	"math"
	"math/rand"
	"time"
)

//var ExtraDataLength int = 0  //int value is zero,



func (c *core) sendDSelect() {
	logger := c.logger.New("state", c.state)
	var extra [7]common.ValidatorInfo  // 최소 7개 노드에서 추가?  기동시 최소 7개 이상 띄워야함.
	//var extra [50]ValidatorInfo  //debugging... 7 -> 50 ,, logical bug. 임시로 50개,, 나중에 수정할 것.
	flag := false

	for i, v := range c.valSet.List() {
		validatorInfo := common.ValidatorInfo{}
		validatorInfo.Address = v.Address()
		validatorInfo.Qrnd = rand.Uint64()

		if i == 0 {
			if !c.valSet.IsProposer(v.Address()) {
				validatorInfo.Tag = common.Coordinator
			} else {
				flag = true
			}
		} else if i == 1 {
			if flag {
				validatorInfo.Tag = common.Coordinator
			}
		} else {
			validatorInfo.Tag = common.Candidate
		}

		extra[i] = validatorInfo
	}

	extraDataJson, err := json.Marshal(extra)
	if err != nil {
		logger.Error("Failed to encode JSON", err)
	}

	c.broadcast(&message{
		Code: msgDSelect,
		Msg: extraDataJson,
	})
}

func (c *core) sendExtraDataRequest() {
	//c.send(&message{
	//	Code: msgExtraDataRequest,
	//	Msg: []byte("extra data request testing."),
	//}, c.qmanager)

	log.Info("Extra Data Requesting To Standalone QMANAGER")

	extra := podc_global.RequestExtraData(c.valSet.GetProposer().Address().String())
	log.Info("Recieved ExtraData ", extra )

	extraDataJson, err := json.Marshal(extra)
	log.Info("Recieved extraDataJson ", extraDataJson )

	if err != nil {
		log.Error("Failed to encode", "extra data", extra)
	}
	c.ExtraDataLength = len(extra)


	c.broadcast(&message{
		Code: msgDSelect,
		Msg: extraDataJson,
	})
}

func (c *core) sendCoordinatorDecide() {
	coordinatorData := c.valSet.GetProposer()
	encodedCoordinatorData, err := Encode(&coordinatorData)

	if err != nil {
		log.Error("Failed to encode", "extra data", coordinatorData)
		return
	}

	c.multicast(&message{
		Code: msgCoordinatorDecide,
		Msg: encodedCoordinatorData,
	}, c.GetValidatorListExceptQman())
}

func (c *core) sendRacing(addr common.Address) {
	c.send(&message{
		Code: msgRacing,
		Msg: []byte("racing testing"),
	}, addr)
}


func (c *core) sendCandidateDecide() {
	c.multicast(&message{
		Code: msgCandidateDecide,
		Msg: []byte("Candidate decide testing"),
	}, c.GetValidatorListExceptQman())
}
//D-Select msg
func (c *core) handleSentExtraData(msg *message, src podc.Validator) error {
	// Decode d-select message
	var extraData common.ValidatorInfos
	if err := json.Unmarshal(msg.Msg, &extraData); err != nil {
		log.Error("JSON Decode Error", "Err", err)
		log.Info("Decode Error")
		return errFailedDecodePrepare
	}
	c.ExtraDataLength = len(extraData)


	c.broadcast(&message{
		Code: msgDSelect,
		Msg: msg.Msg,
	})

	return nil
}

func (c *core) handleDSelect(msg *message, src podc.Validator) error {
	log.Info("4. Get extra data and start d-select", "elapsed", common.PrettyDuration(time.Since(c.intervalTime)))
	c.racingFlag = false
	c.count = 0
	c.intervalTime = time.Now()

	// Decode d-select message
	var extraData common.ValidatorInfos

	if err := json.Unmarshal(msg.Msg, &extraData); err != nil {
		log.Error("JSON Decode Error", "Err", err)
		log.Info("Decode Error")
		return errFailedDecodePrepare
	}
	nodename, err := os.Getwd()
	if err != nil {
		fmt.Printf("current nodename= %v , err=%v",  nodename, err)
	}

	var QRND uint64
	for _, v := range extraData {
		if v.Address == c.address {
			c.tag = v.Tag
			QRND = v.Qrnd
		}
	}


	if c.tag == common.Coordinator {

		//QRNDArray := make([]byte, 8)
		//binary.LittleEndian.PutUint64(QRNDArray, QRND)




		c.ExtraDataLength = len(extraData)
		c.criteria = 29

		isCoordinator := podc_global.CooridnatorConfirmation(podc_global.RequestCoordiStruct{QRND: QRND})

		if isCoordinator{

			var err error
			log.Info(fmt.Sprintf("I am Coordinator! ExtraDataLength %d", c.ExtraDataLength))  //grep -r 'I am Coordinator!' *.log
			if( c.ExtraDataLength != 0 ){
				c.criteria = math.Ceil((( float64(c.ExtraDataLength) - 1.00) * float64(0.51) ) )  //Ceil.. >= 수 리턴.
			}
			if( c.ExtraDataLength == 0 ){
				log.Info("ExtraDataLength has problem")
				//utils.Fatalf("ExtraDataLength has problem)
				return err
			}


			log.Info("c.criteria=", "c.criteria", c.criteria )
			c.sendCoordinatorDecide()
		}


	} else{
		c.ExtraDataLength = 0
		c.criteria = 0
	}

	return nil
}

func (c *core) handleCoordinatorConfirm(msg *message, src podc.Validator) error {
        var err error
		log.Info(fmt.Sprintf("I am Coordinator! ExtraDataLength %d", c.ExtraDataLength))  //grep -r 'I am Coordinator!' *.log
		if( c.ExtraDataLength != 0 ){
			c.criteria = math.Ceil((( float64(c.ExtraDataLength) - 1.00) * float64(0.51) ) )  //Ceil.. >= 수 리턴.
		}
		if( c.ExtraDataLength == 0 ){
			log.Info("ExtraDataLength has problem")
			//utils.Fatalf("ExtraDataLength has problem)
			return err
		}


		log.Info("c.criteria=", "c.criteria", c.criteria )
		c.sendCoordinatorDecide()

	return nil
}

func (c *core) handleCoordinatorDecide(msg *message, src podc.Validator) error {
	if c.tag != common.Coordinator {
		c.sendRacing(src.Address())  //레이싱 시작 메시지 전송
	}

	return nil
}

func (c *core) handleRacing(msg *message, src podc.Validator) error {
	c.racingMu.Lock()
	defer c.racingMu.Unlock()
	if c.tag == common.Coordinator {

		c.count = c.count + 1
		if c.count > uint(c.criteria) && !c.racingFlag {
			c.racingFlag = true
			c.sendCandidateDecide()
		}
	}

	return nil
}

func (c *core) handleCandidateDecide(msg *message, src podc.Validator) error {  //커밋단계로 진입
	if c.state == StatePreprepared {
		log.Info("5. Racing complete and d-select finished.", "elapsed", common.PrettyDuration(time.Since(c.intervalTime)))
		c.intervalTime = time.Now()
		c.setState(StateDSelected)  //D-selected 상태로 설정하고, 커밋 상태로 진입.
		c.sendDCommit()  // msgCommit 를 통하여, 메시지핸들러에서, handleDCommit를 실행, 여기서 c.verifyDCommit에서 inconsistent 발생,
	}

	return nil
}
