package security

import (
	"net"
	"time"
)

type FirewallRule struct {
	Name      string
	Action    RuleAction
	IPRange   *net.IPNet
	Ports     []int
	Protocol  string
	StartTime time.Time
	EndTime   time.Time
}

type RuleAction int

const (
	Allow RuleAction = iota
	Deny
	Log
)

func (fw *Firewall) AddRule(rule FirewallRule) {
	fw.rules = append(fw.rules, rule)
}

func (fw *Firewall) CheckIP(ip net.IP) bool {
	if fw.blacklist[ip.String()] {
		return false
	}

	for _, rule := range fw.rules {
		if rule.IPRange.Contains(ip) {
			return rule.Action == Allow
		}
	}

	return true
}
