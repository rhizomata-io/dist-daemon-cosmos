package inner

import (
	"fmt"
	"log"
	
	"github.com/tendermint/tendermint/types/time"
	
	"github.com/rhizomata-io/dist-daemon-cosmos/node"
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/types"
)

func BroadcastHeartbeat() {
	nodeID := node.GetNodeID()
	heartbeat := types.NewMsgHeartbeat(nodeID,time.Now(), node.GetOpAddress())
	err := BroadcastTx(heartbeat)
	if err != nil {
		log.Println("[ERROR] BroadcastHeartbeat::", err)
	} else {
		heartbeat := types.MsgHeartbeat{}
		Query(heartbeat)
		fmt.Println("----- Heartbeat::",heartbeat)
	}
}
