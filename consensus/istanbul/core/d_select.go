package core

import (
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/log"
	"math/rand"
	"encoding/json"
)

const criteria = 2

func (c *core) sendDSelect() {
	logger := c.logger.New("state", c.state)
	var extra [10]istanbul.Validator

	for i, v := range c.valSet.List() {
		validatorInfo := v
		validatorInfo.SetQrnd(rand.Uint64())

		if c.valSet.IsProposer(v.Address()) {
			validatorInfo.SetTag(istanbul.Coordinator)
		} else {
			validatorInfo.SetTag(istanbul.Coordinator)
		}

		extra[i] = validatorInfo
	}

	extraDataJson, err := json.Marshal(extra)
	if err != nil {
		logger.Error("Failed to encode JSON", err)
	}

	//encodedExtraData, err := Encode(&extra)
	//logger.Info("encode", "extra data", encodedExtraData)
	//if err != nil {
	//	logger.Error("Failed to encode", "extra data", extra)
	//	return
	//}

	c.broadcast(&message{
		Code: msgDSelect,
		Msg:  []byte(extraDataJson),
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

func (c *core) sendRacing() {
	log.Info("sending racing")
	c.send(&message{
		Code: msgRacing,
		Msg: []byte("racing testing"),
	}, c.valSet.GetProposer().Address())
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
	//var extraData *[10]istanbul.Validator
	//
	//err := msg.Decode(&extraData)
	//if err != nil {
	//	return errFailedDecodePrepare
	//}

	//if c.state == StatePreprepared {
	//	c.sendPrepare()
	//}

	if c.valSet.IsProposer(c.Address()) {
		log.Info("I'm proposer!!!!!!")
		c.sendCoordinatorDecide()
	}

	return nil
}

func (c *core) handleCoordinatorDecide(msg *message, src istanbul.Validator) error {
	if !c.valSet.IsProposer(c.Address()) {
		log.Info("handling coordinator decide")
		c.sendRacing()
	}

	return nil
}

func (c *core) handleRacing(msg *message, src istanbul.Validator) error {
	if c.valSet.IsProposer(c.Address()) {
		log.Info("handling racing")
		c.count = c.count + 1

		if c.count > criteria {
			log.Info("send Candidate Decide")
			c.sendCandidateDecide()
		}
	}

	return nil
}

func (c *core) handleCandidateDecide(msg *message, src istanbul.Validator) error {
	if c.state == StatePreprepared {
		log.Info("send prepare!")
		c.sendPrepare()
	}

	return nil
}
