package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
)

// NameKV struct
type NameKV struct {
	Key   string   `json:"key"`
	Value eos.Name `json:"value"`
}

// StringKV struct
type StringKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AssetKV struct
type AssetKV struct {
	Key   string    `json:"key"`
	Value eos.Asset `json:"value"`
}

// TimePointKV struct
type TimePointKV struct {
	Key   string             `json:"key"`
	Value eos.BlockTimestamp `json:"value"`
}

// IntKV struct
type IntKV struct {
	Key   string `json:"key"`
	Value uint64 `json:"value"`
}

// TrxKV struct
type TrxKV struct {
	Key   string          `json:"key"`
	Value eos.Transaction `json:"value"`
}

// FloatKV struct
type FloatKV struct {
	Key   string       `json:"key"`
	Value eos.Float128 `json:"value"`
}

type Scope struct {
	Code  eos.Name `json:"code"`
	Scope eos.Name `json:"scope"`
	Table eos.Name `json:"table"`
	Payer eos.Name `json:"payer"`
	Count uint64   `json:"count"`
}

// Object struct
type Object struct {
	ID           uint64             `json:"id"`
	Names        []NameKV           `json:"names"`
	Strings      []StringKV         `json:"strings"`
	Assets       []AssetKV          `json:"assets"`
	TimePoints   []TimePointKV      `json:"time_points"`
	Ints         []IntKV            `json:"ints"`
	Transactions []TrxKV            `json:"trxs"`
	Floats       []FloatKV          `json:"floats"`
	CreatedDate  eos.BlockTimestamp `json:"created_date"`
	UpdatedDate  eos.BlockTimestamp `json:"updated_date"`
}

// LoadObjects from a provided eos.API and scope
func LoadObjects(ctx context.Context, api *eos.API, scope string) []Object {
	var objects []Object
	var request eos.GetTableRowsRequest
	request.Code = "dao.hypha"
	request.Scope = scope
	request.Table = "objects"
	request.Limit = 1000
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&objects)
	return objects
}

// DAOObject is a generic object from the objects table
type DAOObject struct {
	ID           uint64                        `json:"id"`
	Names        map[string]eos.Name           `json:"names"`
	Strings      map[string]string             `json:"strings"`
	Assets       map[string]eos.Asset          `json:"assets"`
	TimePoints   map[string]eos.BlockTimestamp `json:"time_points"`
	Ints         map[string]uint64             `json:"ints"`
	Transactions map[string]eos.Transaction    `json:"trxs"`
	Floats       map[string]eos.Float128       `json:"floats"`
	CreatedDate  eos.BlockTimestamp            `json:"created_date"`
	UpdatedDate  eos.BlockTimestamp            `json:"updated_date"`
}

// ToDAOObject Converts a very generic Object to one of type DAO Object
func ToDAOObject(objs Object) DAOObject {

	var daoObject DAOObject
	daoObject.Names = make(map[string]eos.Name)
	for index, element := range objs.Names {
		daoObject.Names[element.Key] = objs.Names[index].Value
	}

	daoObject.Assets = make(map[string]eos.Asset)
	for index, element := range objs.Assets {
		daoObject.Assets[element.Key] = objs.Assets[index].Value
	}

	daoObject.TimePoints = make(map[string]eos.BlockTimestamp)
	for index, element := range objs.TimePoints {
		daoObject.TimePoints[element.Key] = objs.TimePoints[index].Value
	}

	daoObject.Ints = make(map[string]uint64)
	for index, element := range objs.Ints {
		daoObject.Ints[element.Key] = objs.Ints[index].Value
	}

	daoObject.Transactions = make(map[string]eos.Transaction)
	for index, element := range objs.Transactions {
		daoObject.Transactions[element.Key] = objs.Transactions[index].Value
	}

	daoObject.Strings = make(map[string]string)
	for index, element := range objs.Strings {
		daoObject.Strings[element.Key] = objs.Strings[index].Value
	}
	daoObject.ID = objs.ID
	daoObject.CreatedDate = objs.CreatedDate
	daoObject.UpdatedDate = objs.UpdatedDate
	return daoObject
}
