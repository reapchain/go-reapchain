package core

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/log"
	"math/rand"
	"time"
)

const criteria = 2

var startTime time.Time
var endTime time.Time

var tstartTime time.Time
var tendTime time.Time

type ValidatorInfo struct {
	Address common.Address
	Tag istanbul.Tag
	Qrnd uint64
}

type ValidatorInfos []ValidatorInfo

func (c *core) sendDSelect() {
	logger := c.logger.New("state", c.state)
	var extra [6]ValidatorInfo
	flag := false

	for i, v := range c.valSet.List() {
		validatorInfo := ValidatorInfo{}
		validatorInfo.Address = v.Address()
		validatorInfo.Qrnd = rand.Uint64()

		if i == 0 {
			if !c.valSet.IsProposer(v.Address()) {
				validatorInfo.Tag = istanbul.Coordinator
			} else {
				flag = true
			}
		} else if i == 1 {
			if flag {
				validatorInfo.Tag = istanbul.Coordinator
			}
		} else {
			validatorInfo.Tag = istanbul.Candidate
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

func (c *core) handleDSelect(msg *message, src istanbul.Validator) error {
	log.Info("d-select start")
	tstartTime = time.Now()
	// Decode d-select message
	var extraData ValidatorInfos

	if err := json.Unmarshal(msg.Msg, &extraData); err != nil {
		log.Error("JSON Decode Error", "Err", err)
		return errFailedDecodePrepare
	}

	for _, v := range extraData {
		if v.Address == c.address {
			c.tag = v.Tag
		}
	}

	if c.tag == istanbul.Coordinator {
		log.Info("I am Coordinator!")
		startTime = time.Now()
		c.sendCoordinatorDecide()
	}

	return nil
}

func (c *core) handleCoordinatorDecide(msg *message, src istanbul.Validator) error {
	if c.tag != istanbul.Coordinator {
		c.sendRacing(src.Address())
	}

	return nil
}

func (c *core) handleRacing(msg *message, src istanbul.Validator) error {
	if c.tag == istanbul.Coordinator {
		c.count = c.count + 1
		log.Info("handling racing", "count", c.count)

		if c.count > criteria {
			endTime = time.Now()
			elapse := (endTime.UnixNano() - startTime.UnixNano()) / 1000000
			log.Info("racing complete", "elapse time(ms)", elapse)
			c.sendCandidateDecide()
			c.count = 0
		}
	}

	return nil
}

func (c *core) handleCandidateDecide(msg *message, src istanbul.Validator) error {
	if c.state == StatePreprepared {
		tendTime = time.Now()
		telapse := (tendTime.UnixNano() - tstartTime.UnixNano()) / 1000000
		log.Info("d-select end", "elapse time(ms)", telapse)
		c.sendDCommit()
	}

	return nil
}
