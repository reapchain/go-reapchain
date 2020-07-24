// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// PoDC is compatible with eth63 protocol  //?
	PoDCName           = "podc"
	PoDCVersion        = 64
	PoDCProtocolLength = 18
	PoDCMsg = 0x11
	QmanMsg = 0x12
)

type PoDCProtocolManager struct {
	*protocolManager

	engine   consensus.PoDC
	eventSub *event.TypeMuxSubscription
}

func newPoDCProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkId uint64, maxPeers int, mux *event.TypeMux, txpool txPool, engine consensus.PoDC, blockchain *core.BlockChain, chaindb ethdb.Database) (*PoDCProtocolManager, error) {
	// Create eth63 protocol manager
	defaultManager, err := newProtocolManager(config, mode, networkId, maxPeers, mux, txpool, engine, blockchain, chaindb)
	if err != nil {
		return nil, err
	}

	// Create the PoDC protocol manager
	manager := &PoDCProtocolManager{
		protocolManager: defaultManager,
		engine:          engine,
	}

	// Support only Istanbul protocol
	manager.SubProtocols = []p2p.Protocol{
		p2p.Protocol{
			Name:    PoDCName,
			Version: PoDCVersion,
			Length:  PoDCProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := manager.newPeer(int(PoDCVersion), p, rw)
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return manager.handle(peer, manager.handleMsg)
//
				case manager.ValidatorSyncCh <- peer:
					manager.wg.Add(1)  //? 양수, 음수 구분 및 나중에,, 추가 검
					defer manager.wg.Done()
					return manager.handle(peer, manager.handleMsg)

					//

				case <-manager.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return manager.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := manager.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		},
	}

	return manager, nil
}
// 비로소 이스탄불 프로토콜 매니저가 시작되는 부분 중
func (pm *PoDCProtocolManager) Start(qman []*discover.Node) {
	// Subscribe required events
	pm.eventSub = pm.eventMux.Subscribe(podc.ConsensusDataEvent{}, core.ChainHeadEvent{}) //이벤트 구독 등록
	//Qmanager로부터 오는 이벤트도 등록해야하나?
	//이스탄불데이터이벤트에 일단은 포함시킨다. Qmanager와 데이터교환을 ConsensusDataEvent에 부분으로 등록한다.

	//Qmanager list에서 하나의 address만 뽑는다.
	go pm.eventLoop()  // //고루틴으로 동시성 처리 // 순서 중요할듯 1. 이벤트루프
	pm.protocolManager.Start(qman) //    2. 프로토콜매니저의 일반 핸들러 시작
	pm.engine.Start(pm.protocolManager.blockchain, qman, pm.commitBlock)  // 합의 엔진 핸들러 시작
}

func (pm *PoDCProtocolManager) Stop() {
	log.Info("Stopping Ethereum protocol")
	pm.engine.Stop()
	pm.protocolManager.Stop()
	pm.eventSub.Unsubscribe() // quits eventLoop
}

// handleMsg handles PoDC related consensus messages or
// fallback to default procotol manager's handler
// PoDC handler and  general handler decision and conditional jump


func (pm *PoDCProtocolManager) handleMsg(p *peer, msg p2p.Msg) error {
	// Handle Istanbul messages
	switch {
	case msg.Code == PoDCMsg:
		pubKey, err := p.ID().Pubkey()  //enode://"public key@IP address: port number"
		                                //if PoDC message, for example a response ( extradata from Qmanager)
		                                // extract pubkey and err
		if err != nil {
			return err
		}
		var data []byte
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		return pm.engine.HandleMsg(pubKey, data)  // -> PoDC message handler : type PoDC interface.. /consensus.go
	//









	//

	default:
		// Invoke default protocol manager's message handler
		return pm.protocolManager.handleMsg(p, msg) // -> general protocol message handler /eth/handler.go
	}
}

// event loop for PoDC
func (pm *PoDCProtocolManager) eventLoop() {
	for obj := range pm.eventSub.Chan() {
		switch ev := obj.Data.(type) {
		case podc.ConsensusDataEvent:
			pm.sendEvent(ev)
		case core.ChainHeadEvent:
			pm.newHead(ev)
		}
	}
}

// sendEvent sends a p2p message with given data to a peer
func (pm *PoDCProtocolManager) sendEvent(event podc.ConsensusDataEvent) {
	// FIXME: it's inefficient because it retrieves all peers every time
	p := pm.findPeer(event.Target)
	if p == nil {
		log.Warn("Failed to find peer by address", "addr", event.Target)
		return
	}
	p2p.Send(p.rw, PoDCMsg, event.Data)
}

func (pm *PoDCProtocolManager) commitBlock(block *types.Block) error {
	if _, err := pm.blockchain.InsertChain(types.Blocks{block}); err != nil {
		log.Debug("Failed to insert block", "number", block.Number(), "hash", block.Hash(), "err", err)
		return err
	}
	// Only announce the block, don't broadcast it
	go pm.BroadcastBlock(block, false)
	return nil
}

func (pm *PoDCProtocolManager) newHead(event core.ChainHeadEvent) {
	block := event.Block
	if block != nil {
		pm.engine.NewChainHead(block)
	}
}

// findPeer retrieves a peer by given address
func (pm *PoDCProtocolManager) findPeer(addr common.Address) *peer {
	for _, p := range pm.peers.Peers() {
		pubKey, err := p.ID().Pubkey()
		if err != nil {
			continue
		}
		if crypto.PubkeyToAddress(*pubKey) == addr {
			fmt.Sprintf("%x", addr)
			return p
		}
	}
	return nil
}
