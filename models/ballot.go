package models

import (
	"context"
	"errors"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

// Ballot ...
type Ballot struct {
	BallotName  eos.Name             `json:"ballot_name"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Status      string               `json:"status"`
	VoteTally   map[string]eos.Asset `json:"options"`
	Votes       []Vote
	PassVotes   eos.Asset
	RejectVotes eos.Asset
	BeginTime eos.BlockTimestamp `json:"begin_time"`
	EndTime   eos.BlockTimestamp `json:"end_time"`
}

// Vote ...
type Vote struct {
	Voter          eos.Name           `json:"voter"`
	WeightedVotes  []AssetKV          `json:"weighted_votes"`
	VotingPower    eos.Asset          `json:"raw_votes"`
	VotingTime     eos.BlockTimestamp `json:"vote_time"`
	VoteSelections map[string]eos.Asset
}

//NewBallot ...
func NewBallot(ctx context.Context, api *eos.API, ballotName eos.Name) (*Ballot, error) {
	var ballot []Ballot
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("TelosDecideContract")
	request.Scope = viper.GetString("TelosDecideContract")
	request.Table = "ballots"
	request.Limit = 1
	request.LowerBound = string(ballotName)
	request.UpperBound = string(ballotName)
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return nil, err
	}
	response.JSONToStructs(&ballot)
	if len(ballot) == 0 {
		return nil, errors.New("Ballot name was not found: " + string(ballotName))
	}

	var votes []Vote
	var voteRequest eos.GetTableRowsRequest
	voteRequest.Code = viper.GetString("TelosDecideContract")
	voteRequest.Scope = string(ballot[0].BallotName)
	voteRequest.Table = "votes"
	voteRequest.Limit = 500
	voteRequest.JSON = true
	voteResponse, err := api.GetTableRows(ctx, voteRequest)
	if err != nil {
		return nil, err
	}
	voteResponse.JSONToStructs(&votes)

  votesAgainstTotal, _ := eos.NewAssetFromString("0.00 HVOICE")
  votesForTotal, _ := eos.NewAssetFromString("0.00 HVOICE")

	for _, vote := range votes {
		vote.VoteSelections = make(map[string]eos.Asset)
		for index, selection := range vote.WeightedVotes {
			vote.VoteSelections[selection.Key] = vote.WeightedVotes[index].Value
			if selection.Key == "pass" {
			  votesForTotal = votesForTotal.Add(vote.WeightedVotes[index].Value)
      } else {
        votesAgainstTotal = votesAgainstTotal.Add(vote.WeightedVotes[index].Value)
      }
		}
	}

  ballot[0].PassVotes = votesForTotal
  ballot[0].RejectVotes = votesAgainstTotal
	ballot[0].Votes = votes
	return &ballot[0], nil
}

// GetHvoiceSupply ...
func GetHvoiceSupply(ctx context.Context, api *eos.API) (*eos.Asset, error) {
	type Supply struct {
		HvoiceSupply eos.Asset `json:"supply"`
	}

	var supply []Supply
	// telosDecide := eos.MustStringToName(viper.GetString("TelosDecideContract"))

	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("TelosDecideContract")
	request.Scope = viper.GetString("TelosDecideContract")
	request.Table = "treasuries"
	request.Limit = 1
	request.LowerBound = string("HVOICE")
	request.UpperBound = string("HVOICE")
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return nil, err
	}
	response.JSONToStructs(&supply)
	return &supply[0].HvoiceSupply, nil
}
