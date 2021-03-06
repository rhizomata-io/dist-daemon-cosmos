package types

import (
	"fmt"
	"strings"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)



// Member member info
type Member struct {
	NodeID     string `json:"nodeid"`
	Name       string `json:"name"`
	Heartbeat  time.Time `json:"heartbeat"`
	Owner      sdk.AccAddress `json:"owner"`
	leader     bool // transient field
	alive      bool // transient field
	local      bool // transient field
}

var (
	NilMember = Member{}
)

func NewMember(name string, nodeid string, owner sdk.AccAddress) Member{
	return Member{NodeID:nodeid, Name:name, Owner:owner}
}

// IsLeader return whether member is leader
func (m *Member) IsLeader() bool {
	return m.leader
}

// SetLeader Set member as leader
func (m *Member) SetLeader(leader bool) {
	m.leader = leader
}

// IsAlive return whether member is alive
func (m *Member) IsAlive() bool {
	return m.alive
}

// SetAlive Set member alive
func (m *Member) SetAlive(alive bool) {
	m.alive = alive
}

// IsLocal return whether member is alive
func (m *Member) IsLocal() bool {
	return m.local
}

// SetLocal Set member alive
func (m *Member) SetLocal(local bool) {
	m.local = local
}


// String implement fmt.Stringer
func (m Member) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Member[%s:%s] alive %t`,
		m.Name, m.NodeID, m.alive))
}
