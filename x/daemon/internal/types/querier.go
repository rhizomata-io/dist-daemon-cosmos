package types

import (
	"strings"
	"fmt"
)

// QueryMember Queries Result Payload for a member
type QueryResMember struct {
	NodeID string `json:"nodeid"`
	Name string `json:"name"`
}

// implement fmt.Stringer
func (r QueryResMember) String() string {
	return fmt.Sprintf("%s:%s",r.Name, r.NodeID)
}

// QueryMembers Queries Result Payload for member ids query
type QueryResNodeIDs []string

// implement fmt.Stringer
func (n QueryResNodeIDs) String() string {
	return strings.Join(n[:], "\n")
}
