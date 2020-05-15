package models

import eos "github.com/eoscanada/eos-go"

// Member ...
type Member struct {
	MemberName        eos.Name
	CurrentAssignment Assignment
}
