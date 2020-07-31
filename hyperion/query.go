package hyperion

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hypha-dao/daoctl/models"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// Query represents an object to pass off to Hyperion to generate results
type Query struct {
	Action   string
	Contract string
	Account  string
	TrxID    string
	After    time.Time
	Before   time.Time
	Limit    int
}

// NewQuery returns a query object that can be updated and then executed
func NewQuery(action, contract, account string) Query {
	return Query{
		Action:   action,
		Contract: contract,
		Account:  account,
		Limit:    10, // defaults to 10
	}
}

// Results returns the query results from Hyperiond
func (q *Query) Results() ([]models.QrAction, error) {
	request := viper.GetString("HyperionEndpoint") + "/history/"
	if q.TrxID != "" {
		request += "get_transaction?id=" + q.TrxID
	} else {
		request += "get_actions?limit=" + strconv.Itoa(q.Limit)
		if q.Action != "" {
			request += "&act.name=" + q.Action
		}
		if q.Contract != "" {
			request += "&act.account=" + q.Contract
		}
		if q.Account != "" {
			request += "&account=" + q.Account
		}
		if !q.After.IsZero() {
			request += "&after=" + q.After.Format("2006-01-02T15:04:05.000")
		}
		if !q.Before.IsZero() {
			request += "&before=" + q.After.Format("2006-01-02T15:04:05.000")
		}
	}

	resp, err := http.Get(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var actions []models.QrAction
	result := gjson.Get(string(body), "actions")
	result.ForEach(func(key, value gjson.Result) bool {
		actionTime, _ := time.Parse("2006-01-02T15:04:05.000", getString(value, "timestamp"))

		action := models.QrAction{
			Timestamp:      actionTime,
			TrxID:          gjson.Get(value.String(), "trx_id").String(),
			ActionContract: getString(value, "act.account"),
			ActionName:     getString(value, "act.name"),
			Data:           getString(value, "act.data"),
		}
		actions = append(actions, action)
		return true // keep iterating
	})
	return actions, nil
}

func getString(result gjson.Result, element string) string {

	charsToShow := 45
	suffix := "... <snip>"
	longValue := gjson.Get(result.String(), element).String()

	if len(longValue) < charsToShow {
		charsToShow = len(longValue)
		suffix = ""
	}
	return longValue[:charsToShow] + suffix
}
