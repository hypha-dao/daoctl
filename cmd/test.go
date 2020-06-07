package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/spf13/cobra"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"12345"

const creator = "eosio"
const defaultKey = ""
const repoHome = "/Users/max/dev/hypha/eosio-contracts"

var testingKey ecc.PublicKey

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randAccountName() string {
	return stringWithCharset(12, charset)
}

func execTrx(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		log.Printf("Error filling tx opts: %s", err)
		return "error", err
	}

	tx := eos.NewTransaction(actions, txOpts)
	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		log.Printf("Error signing transaction: %s", err)
		return "error", err
	}

	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		log.Printf("Error pushing transaction: %s", err)
		return "error", err
	}
	trxID := hex.EncodeToString(response.Processed.ID)
	return trxID, nil
}

type tDTreasury struct {
	Manager   eos.AccountName
	MaxSupply eos.Asset
	Access    eos.Name
}

func newTreasury(ctx context.Context, api *eos.API, telosDecide, treasuryManager *eos.AccountName) (string, error) {
	maxSupply, _ := eos.NewAssetFromString("1000000000.00 HVOICE")
	actions := []*eos.Action{
		{
			Account: *telosDecide,
			Name:    toActionName("newtreasury", "creating new treasury within Telos Decide"),
			Authorization: []eos.PermissionLevel{
				{Actor: *treasuryManager, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(tDTreasury{
				Manager:   *treasuryManager,
				MaxSupply: maxSupply,
				Access:    eos.Name("public"),
			}),
		}}
	return execTrx(ctx, api, actions)
}

func createRandoms(ctx context.Context, api *eos.API, length int) ([]*eos.AccountName, error) {

	i := 0
	var actions []*eos.Action
	var accounts []*eos.AccountName
	accounts = make([]*eos.AccountName, length)
	keyBag := api.Signer

	for i < length {
		acct := toAccount(randAccountName(), "random account name")
		key, _ := ecc.NewRandomPrivateKey()

		err := keyBag.ImportPrivateKey(ctx, key.String())
		if err != nil {
			log.Panicf("import private key: %s", err)
		}

		accounts[i] = &acct
		actions = append(actions, system.NewNewAccount(creator, acct, key.PublicKey()))
		log.Println("Creating account: 	", acct, " with private key : ", key.String())
		i++
	}

	trxID, err := execTrx(ctx, api, actions)
	if err != nil {
		log.Panicf("cannot create random accounts: %s", err)
		return nil, err
	}
	log.Println("Created random accounts: ", trxID)
	return accounts, nil
}

func setContract(ctx context.Context, api *eos.API, accountName *eos.AccountName, wasmFile, abiFile string) (string, error) {
	setCodeAction, err := system.NewSetCode(*accountName, wasmFile)
	errorCheck("loading wasm file", err)

	setAbiAction, err := system.NewSetABI(*accountName, abiFile)
	errorCheck("loading abi file", err)

	return execTrx(ctx, api, []*eos.Action{setCodeAction, setAbiAction})
}

func setConfig(ctx context.Context, api *eos.API, contract *eos.AccountName, configFile string) (string, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("cannot read configuration: %s", err)
		return "error", err
	}

	action := toActionName("setconfig", "action")

	var dump map[string]interface{}
	err = json.Unmarshal(data, &dump)
	if err != nil {
		log.Panicf("cannot unmarshal configuration: %s", err)
		return "error", err
	}
	actionBinary, err := api.ABIJSONToBin(ctx, *contract, eos.Name(action), dump)
	errorCheck("unable to retrieve action binary from JSON via API", err)

	actions := []*eos.Action{
		{
			Account: *contract,
			Name:    action,
			Authorization: []eos.PermissionLevel{
				{Actor: *contract, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionDataFromHexData([]byte(actionBinary)),
		}}

	return execTrx(ctx, api, actions)
}

type appVersion struct {
	AppVersion string
}

func initTD(ctx context.Context, api *eos.API, telosDecide eos.AccountName) (string, error) {
	actions := []*eos.Action{
		{
			Account: telosDecide,
			Name:    toActionName("init", "init action name on Telos Decide"),
			Authorization: []eos.PermissionLevel{
				{Actor: telosDecide, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(appVersion{
				AppVersion: "vtest",
			}),
		}}
	return execTrx(ctx, api, actions)
}

type fee struct {
	FeeName   eos.Name
	FeeAmount eos.Asset
}

func setFee(ctx context.Context, api *eos.API, telosDecide eos.AccountName) (string, error) {
	zeroFee, _ := eos.NewAssetFromString("0.0000 TLOS")
	feeNames := []eos.Name{eos.Name("ballot"), eos.Name("treasury"), eos.Name("archival"), eos.Name("committee")}
	var fees []*eos.Action
	for _, feeName := range feeNames {
		fee := eos.Action{
			Account: telosDecide,
			Name:    toActionName("updatefee", "td update fee action"),
			Authorization: []eos.PermissionLevel{
				{Actor: telosDecide, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(fee{
				FeeName:   feeName,
				FeeAmount: zeroFee,
			}),
		}
		fees = append(fees, &fee)
	}
	return execTrx(ctx, api, fees)
}

type addPeriod struct {
	StartTime eos.TimePoint `json:"start_time"`
	EndTime   eos.TimePoint `json:"end_time"`
	Phase     string        `json:"phase"`
}

func addPeriods(ctx context.Context, api *eos.API, daoContract eos.AccountName, numPeriods int, periodDuration time.Duration) (string, error) {

	now := time.Now()

	startTime := eos.TimePoint(now.UnixNano() / 1000)
	endTime := eos.TimePoint(now.Add(periodDuration).UnixNano() / 1000)

	var periods []*eos.Action

	for i := 0; i < numPeriods; i++ {
		addPeriodAction := eos.Action{
			Account: daoContract,
			Name:    toActionName("addperiod", "add period action name"),
			Authorization: []eos.PermissionLevel{
				{Actor: daoContract, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(addPeriod{
				StartTime: startTime,
				EndTime:   endTime,
				Phase:     "test phase",
			}),
		}
		periods = append(periods, &addPeriodAction)
	}

	return execTrx(ctx, api, periods)
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [id]",
	Short: "Test the DAO software",

	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New("http://localhost:8888")
		// api.Debug = true
		ctx := context.Background()
		wasmFile := repoHome + "/hyphadao/hyphadao.wasm"
		abiFile := repoHome + "/hyphadao/hyphadao.abi"

		keyBag := &eos.KeyBag{}
		err := keyBag.ImportPrivateKey(ctx, defaultKey)
		if err != nil {
			log.Panicf("cannot import default private key: %s", err)
		}
		api.SetSigner(keyBag)

		accounts, err := createRandoms(ctx, api, 10)
		if err != nil {
			log.Panicf("cannot create random accounts: %s", err)
		}

		daoContract := accounts[0]
		trxID, err := setContract(ctx, api, daoContract, wasmFile, abiFile)
		if err != nil {
			log.Panicf("cannot set contract: %s", err)
		}
		log.Println("Set contract: ", trxID)

		trxID, err = setConfig(ctx, api, daoContract, "test-payloads/config1.json")
		if err != nil {
			log.Panicf("cannot set config: %s", err)
		}
		log.Println("Set config: ", trxID)

		telosDecide := accounts[1]
		trxID, err = setContract(ctx, api, telosDecide, "/Users/max/dev/decide/decide/decide.wasm", "/Users/max/dev/decide/decide/decide.abi")
		if err != nil {
			log.Panicf("cannot set contract: %s", err)
		}
		log.Println("Set Telos Decide contract: ", trxID)

		trxID, err = initTD(ctx, api, *telosDecide)
		if err != nil {
			log.Panicf("cannot initialize TD: %s", err)
		}
		log.Println("Initialized telos decide: ", trxID)

		trxID, err = setFee(ctx, api, *telosDecide)
		if err != nil {
			log.Panicf("cannot set Fee on TD: %s", err)
		}
		log.Println("Set fees on TD: ", trxID)

		trxID, err = newTreasury(ctx, api, telosDecide, daoContract)
		if err != nil {
			log.Panicf("cannot create treasury: %s", err)
		}
		log.Println("Created TD treasury: ", trxID)

		fiveMins, _ := time.ParseDuration("5m")
		trxID, err = addPeriods(ctx, api, *daoContract, 10, fiveMins)
		if err != nil {
			log.Panicf("cannot add periods: %s", err)
		}
		log.Println("Added periods to DAO contract: ", trxID)
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
