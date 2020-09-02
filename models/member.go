package models

import (
	"context"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/util"
	"github.com/ryanuber/columnize"
	"github.com/spf13/viper"
)

// Member ...
type Member struct {
	Account            eos.Name
	VoteTokenBalance   eos.Asset
	RewardTokenBalance eos.Asset
}

// TDBalance represents one record in the telos.decide::voters table
type TDBalance struct {
	Liquid eos.Asset `json:"liquid"`
	// only pulling in 'Liquid' balance at this time
}

func (m *Member) String() string {

	output := []string{
		fmt.Sprintf("Member Account|%v", m.Account),
		fmt.Sprintf(viper.GetString("VoteTokenSymbol")+"|%v", util.FormatAsset(&m.VoteTokenBalance, 2)),
		fmt.Sprintf(viper.GetString("RewardToken.Symbol")+"|%v", util.FormatAsset(&m.RewardTokenBalance, 2)),
	}
	return columnize.SimpleFormat(output)
}

// NewMember converts a generic DAO Object to a typed Payout
func NewMember(ctx context.Context, api *eos.API, acct eos.Name) Member {
	var m Member
	var err1 error
	m.Account = acct
	rewardTokenBalance, _ := api.GetCurrencyBalance(ctx, eos.AN(string(acct)), viper.GetString("RewardToken.Symbol"), eos.AN(viper.GetString("RewardTokenContract")))
	if len(rewardTokenBalance) == 0 {
		//fmt.Println("Reward token not found, using 0.00 " + viper.GetString("RewardToken.Symbol"))
		m.RewardTokenBalance, err1 = eos.NewAssetFromString("0.00 " + viper.GetString("RewardToken.Symbol")) // could fail
		if err1 != nil {
			// TODO fix error handling and logging throughout app
			panic("Unable to construct Asset object from the RewardToken.Symbol in configuration; please verify and try again.")
		}
	} else {
		m.RewardTokenBalance = rewardTokenBalance[0]
	}

	var tdb []TDBalance
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("TelosDecideContract")
	request.Scope = string(acct)
	request.Table = "voters"
	request.Limit = 1
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&tdb)

	m.VoteTokenBalance = tdb[0].Liquid // TODO: support users that a members of multiple DAOs
	return m
}

// MemberRecord represents a single row in the dao::members table
type MemberRecord struct {
	MemberName eos.Name `json:"member"`
}

// Members retrieves a list of all of the DAO members, including balances
func Members(ctx context.Context, api *eos.API) []Member {
	var memberRecords []MemberRecord
	// var memberAccounts []eos.Name
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("DAOContract")
	request.Scope = viper.GetString("DAOContract")
	request.Table = "members"
	request.Limit = 1000 // TODO: support dynamic number of members
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&memberRecords)

	var members []Member
	members = make([]Member, len(memberRecords))
	for index, memberRecord := range memberRecords {
		members[index] = NewMember(ctx, api, memberRecord.MemberName)
	}

	return members
}

// ApplicantRecord represents a single row in the dao::members table
type ApplicantRecord struct {
	Applicant eos.Name `json:"applicant"`
}

// Applicants retrieves a list of all of the DAO members, including balances
func Applicants(ctx context.Context, api *eos.API) []ApplicantRecord {
	var applicantRecords []ApplicantRecord
	// var memberAccounts []eos.Name
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("DAOContract")
	request.Scope = viper.GetString("DAOContract")
	request.Table = "applicants"
	request.Limit = 1000 // TODO: support dynamic number of members
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&applicantRecords)
	return applicantRecords
}
