package testing

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eosc/cli"
	"github.com/hypha-dao/daoctl/util"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "12345"
const creator = "eosio"

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

func ExecTrx(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
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

func toName(in, field string) eos.Name {
	name, err := cli.ToName(in)
	if err != nil {
		util.ErrorCheck(fmt.Sprintf("invalid name format for %q", field), err)
	}

	return name
}

// {
//     "account": "dao.gba",
//     "permission": "active",
//     "parent": "owner",
//     "auth": {
//         "keys": [
//             {
//                 "key": "EOS8Y2bVJDB6f1GWLBuyQA73wPSTj4DriHwc1nchS3brQc277BURk",
//                 "weight": 1
//             }
//         ],
//         "threshold": 1,
//         "accounts": [{
// 			"permission": {
// 				"actor": "thomashypha1",
// 				"permission": "active"
// 			},
// 			"weight": 1
// 		}],
//         "waits": []
//     }
// }

func CreateRandoms(ctx context.Context, api *eos.API, length int) ([]*eos.AccountName, error) {

	i := 0
	var actions []*eos.Action
	var accounts []*eos.AccountName
	accounts = make([]*eos.AccountName, length)
	keyBag := api.Signer

	var codePermissionActions []*eos.Action
	codePermissionActions = make([]*eos.Action, length)

	for i < length {
		acct := util.ToAccount(randAccountName(), "random account name")
		key, _ := ecc.NewRandomPrivateKey()

		err := keyBag.ImportPrivateKey(ctx, key.String())
		if err != nil {
			log.Panicf("import private key: %s", err)
		}

		accounts[i] = &acct
		actions = append(actions, system.NewNewAccount(creator, acct, key.PublicKey()))

		codePermissionActions[i] = system.NewUpdateAuth(*accounts[i],
			"active",
			"owner",
			eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{{
					PublicKey: key.PublicKey(),
					Weight:    1,
				}},
				Accounts: []eos.PermissionLevelWeight{{
					Permission: eos.PermissionLevel{
						Actor:      acct,
						Permission: "eosio.code",
					},
					Weight: 1,
				}},
				Waits: []eos.WaitWeight{},
			}, "owner")

		log.Println("Creating account: 	", acct, " with private key : ", key.String())
		i++
	}

	trxID, err := ExecTrx(ctx, api, actions)
	if err != nil {
		log.Panicf("cannot create random accounts: %s", err)
		return nil, err
	}
	log.Println("Created random accounts: ", trxID)

	for _, codePermissionAction := range codePermissionActions {
		trxID, err = ExecTrx(ctx, api, []*eos.Action{codePermissionAction})
		if err != nil {
			log.Panicf("cannot add eosio.code permission: %s", err)
			return nil, err
		}
		log.Println("Added eosio.code permission: ", trxID)
	}

	return accounts, nil
}

func SetContract(ctx context.Context, api *eos.API, accountName *eos.AccountName, wasmFile, abiFile string) (string, error) {
	setCodeAction, err := system.NewSetCode(*accountName, wasmFile)
	util.ErrorCheck("loading wasm file", err)

	setAbiAction, err := system.NewSetABI(*accountName, abiFile)
	util.ErrorCheck("loading abi file", err)

	return ExecTrx(ctx, api, []*eos.Action{setCodeAction, setAbiAction})
}
