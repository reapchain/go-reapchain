package core

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"

	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/log"
	"math"
	"math/rand"
	"time"
)

type ValidatorInfo struct {
	Address common.Address
	Tag podc.Tag
	Qrnd uint64
}

type ValidatorInfos []ValidatorInfo

func (c *core) sendDSelect() {
	logger := c.logger.New("state", c.state)
	var extra [7]ValidatorInfo
	flag := false

	for i, v := range c.valSet.List() {
		validatorInfo := ValidatorInfo{}
		validatorInfo.Address = v.Address()
		validatorInfo.Qrnd = rand.Uint64()

		if i == 0 {
			if !c.valSet.IsProposer(v.Address()) {
				validatorInfo.Tag = podc.Coordinator
			} else {
				flag = true
			}
		} else if i == 1 {
			if flag {
				validatorInfo.Tag = podc.Coordinator
			}
		} else {
			validatorInfo.Tag = podc.Candidate
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
	c.send(&message{
		Code: msgExtraDataRequest,
		Msg: []byte("extra data request testing."),
	}, c.qmanager)
}

func (c *core) sendCoordinatorDecide() {
	coordinatorData := c.valSet.GetProposer()
	encodedCoordinatorData, err := Encode(&coordinatorData)

	if err != nil {
		log.Error("Failed to encode", "extra data", coordinatorData)
		return
	}

	c.broadcast(&message{
		Code: msgCoordinatorDecide,
		Msg: encodedCoordinatorData,
	})
}

func (c *core) sendRacing(addr common.Address) {
	c.send(&message{
		Code: msgRacing,
		Msg: []byte("racing testing"),
	}, addr)
}

func (c *core) sendCandidateDecide() {
	c.broadcast(&message{
		Code: msgCandidateDecide,
		Msg: []byte("Candidate decide testing"),
	})
}

func (c *core) handleSentExtraData(msg *message, src podc.Validator) error {
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
	var extraData ValidatorInfos

	if err := json.Unmarshal(msg.Msg, &extraData); err != nil {
		log.Error("JSON Decode Error", "Err", err)
		log.Info("Decode Error")
		return errFailedDecodePrepare
	}
	log.Info("get extradata display = ", "extraData", extraData)
	nodename, err := os.Getwd()
	if err != nil {
		fmt.Printf("current nodename= %v , err=%v",  nodename, err)
	}
	for _, v := range extraData {
		if v.Address == c.address {
			c.tag = v.Tag
		}
	}

	if c.tag == podc.Coordinator {
		log.Info("I am Coordinator!")
		c.criteria = math.Floor((float64(len(extraData)) - 1) * 0.51)
		log.Info("c.criteria=", "c.criteria", c.criteria )
		c.sendCoordinatorDecide()
	}

	return nil
}

func (c *core) handleCoordinatorDecide(msg *message, src podc.Validator) error {
	if c.tag != podc.Coordinator {
		c.sendRacing(src.Address())
	}

	return nil
}

func (c *core) handleRacing(msg *message, src podc.Validator) error {
	c.racingMu.Lock()
	defer c.racingMu.Unlock()
	if c.tag == podc.Coordinator {

		c.count = c.count + 1
		log.Info("handling racing", "count", c.count)
		log.Info("handling racing", "flag", c.racingFlag)
		if c.count > uint(c.criteria) && !c.racingFlag {
			log.Info("racing completed.", "count", c.count)
			c.racingFlag = true
			c.sendCandidateDecide()
		}
	}

	return nil
}

func (c *core) handleCandidateDecide(msg *message, src podc.Validator) error {
	if c.state == StatePreprepared {
		log.Info("5. Racing complete and d-select finished.", "elapsed", common.PrettyDuration(time.Since(c.intervalTime)))
		c.intervalTime = time.Now()
		c.sendDCommit()
	}

	return nil
}
