package core

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/log"
	"math/rand"
)

const criteria = 2

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
	log.Info("sending coordinator decide")
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
	log.Info("sending racing", "to", addr)
	log.Info("from", "my address", c.address)
	c.send(&message{
		Code: msgRacing,
		Msg: []byte("racing testing"),
	}, addr)
}

func (c *core) sendCandidateDecide() {
	log.Info("sending candidate decide")
	c.broadcast(&message{
		Code: msgCandidateDecide,
		Msg: []byte("Candidate decide testing"),
	})
}

func (c *core) handleDSelect(msg *message, src istanbul.Validator) error {
	// Decode d-select message
	var extraData ValidatorInfos

	if err := json.Unmarshal(msg.Msg, &extraData); err != nil {
		log.Error("JSON Decode Error", "Err", err)
		return errFailedDecodePrepare
	}
	log.Info("JSON Encode", "result", extraData)

	for _, v := range extraData {
		if v.Address == c.address {
			c.tag = v.Tag
		}
	}

	if c.tag == istanbul.Coordinator {
		log.Info("I am Coordinator!!!!")
		c.sendCoordinatorDecide()
	}

	return nil
}

func (c *core) handleCoordinatorDecide(msg *message, src istanbul.Validator) error {
	if c.tag != istanbul.Coordinator {
		log.Info("handling coordinator decide")
		c.sendRacing(src.Address())
	}

	return nil
}

func (c *core) handleRacing(msg *message, src istanbul.Validator) error {
	if c.tag == istanbul.Coordinator {
		c.count = c.count + 1
		log.Info("handling racing", "count", c.count)

		if c.count > criteria {
			log.Info("send Candidate Decide")
			c.sendCandidateDecide()
			c.count = 0
		}
	}

	return nil
}

func (c *core) handleCandidateDecide(msg *message, src istanbul.Validator) error {
	if c.state == StatePreprepared {
		log.Info("send prepare!")
		c.sendDCommit()
	}

	return nil
}
