package cluster

import "nano/cluster/clusterpb"

type Member struct {
	isMaster   bool
	memberInfo *clusterpb.MemberInfo
}

func (m *Member) MemberInfo() *clusterpb.MemberInfo {
	return m.memberInfo
}
