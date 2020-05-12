package models

import eos "github.com/eoscanada/eos-go"

type Member struct {
	MemberName        eos.Name
	CurrentAssignment Assignment
}
