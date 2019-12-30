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
	"github.com/ethereum/go-ethereum/consensus/istanbul"
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
	// istanbul is compatible with eth63 protocol
	istanbulName           = "istanbul"
	istanbulVersion        = 64
	istanbulProtocolLength = 18

	IstanbulMsg = 0x11
	QmanMsg = 0x12 // added by yichoi for qman data exchange
)

type istanbulProtocolManager struct {
	*protocolManager

	engine   consensus.Istanbul
	eventSub *event.TypeMuxSubscription
}

func newIstanbulProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkId uint64, maxPeers int, mux *event.TypeMux, txpool txPool, engine consensus.Istanbul, blockchain *core.BlockChain, chaindb ethdb.Database) (*istanbulProtocolManager, error) {
	// Create eth63 protocol manager
	defaultManager, err := newProtocolManager(config, mode, networkId, maxPeers, mux, txpool, engine, blockchain, chaindb)
	if err != nil {
		return nil, err
	}

	// Create the istanbul protocol manager
	manager := &istanbulProtocolManager{
		protocolManager: defaultManager,
		engine:          engine,
	}

	// Support only Istanbul protocol
	manager.SubProtocols = []p2p.Protocol{
		p2p.Protocol{
			Name:    istanbulName,
			Version: istanbulVersion,
			Length:  istanbulProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := manager.newPeer(int(istanbulVersion), p, rw)
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return manager.handle(peer, manager.handleMsg)
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
// 비로소 이스탄불 프로토콜 매니저가 시작되는 부분 중요
func (pm *istanbulProtocolManager) Start(qman []*discover.Node) {
	// Subscribe required events
	pm.eventSub = pm.eventMux.Subscribe(istanbul.ConsensusDataEvent{}, core.ChainHeadEvent{})  //이벤트 구독 등록
	//Qmanager로부터 오는 이벤트도 등록해야하나?
	//이스탄불데이터이벤트에 일단은 포함시킨다. Qmanager와 데이터교환을 ConsensusDataEvent에 부분으로 등록한다.

	//Qmanager list에서 하나의 address만 뽑는다.

	go pm.eventLoop()  //고루틴으로 동시성 처리
	pm.protocolManager.Start(qman)
	pm.engine.Start(pm.protocolManager.blockchain, qman, pm.commitBlock)  //엔진 시작 : 이스탄불
}

func (pm *istanbulProtocolManager) Stop() {
	log.Info("Stopping Ethereum protocol")
	pm.engine.Stop()
	pm.protocolManager.Stop()
	pm.eventSub.Unsubscribe() // quits eventLoop
}

// handleMsg handles Istanbul related consensus messages or
// fallback to default procotol manager's handler
func (pm *istanbulProtocolManager) handleMsg(p *peer, msg p2p.Msg) error {
	// Handle Istanbul messages
	switch {
	case msg.Code == IstanbulMsg:
		pubKey, err := p.ID().Pubkey()
		if err != nil {
			return err
		}
		var data []byte
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		return pm.engine.HandleMsg(pubKey, data)
	default:
		// Invoke default protocol manager's message handler
		return pm.protocolManager.handleMsg(p, msg)
	}
}

// event loop for Istanbul
func (pm *istanbulProtocolManager) eventLoop() {
	for obj := range pm.eventSub.Chan() {  //채널의 이벤트를 받음.
		switch ev := obj.Data.(type) {
		case istanbul.ConsensusDataEvent:   //여기 내부에서 Qmanager 전송 / 수신 처리할것, 이유는 mux등 복잡한 내부 메카니즘을 쓰기에, 분기 안하는게 좋을듯
			pm.sendEvent(ev)
		//yichoi -Qman data event
		//case istanbul.QmanDataEvent:
		//	pm.sendEvent(ev)
        //end
		case core.ChainHeadEvent:
			pm.newHead(ev)
		}
	}
}

// sendEvent sends a p2p message with given data to a peer
func (pm *istanbulProtocolManager) sendEvent(event istanbul.ConsensusDataEvent) {
	// FIXME: it's inefficient because it retrieves all peers every time
	p := pm.findPeer(event.Target)
	if p == nil {
		log.Warn("Failed to find peer by address", "addr", event.Target)
		return
	}
	p2p.Send(p.rw, IstanbulMsg, event.Data)
}

func (pm *istanbulProtocolManager) commitBlock(block *types.Block) error {
	if _, err := pm.blockchain.InsertChain(types.Blocks{block}); err != nil {  //insert chain
		log.Debug("Failed to insert block", "number", block.Number(), "hash", block.Hash(), "err", err)
		return err
	}
	// Only announce the block, don't broadcast it
	go pm.BroadcastBlock(block, false)
	return nil
}

func (pm *istanbulProtocolManager) newHead(event core.ChainHeadEvent) {
	block := event.Block
	if block != nil {
		pm.engine.NewChainHead(block)
	}
}

// findPeer retrieves a peer by given address
func (pm *istanbulProtocolManager) findPeer(addr common.Address) *peer {
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
