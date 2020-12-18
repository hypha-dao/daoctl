package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	eostest "github.com/digital-scarcity/eos-go-test"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/k0kubun/go-ansi"
	progressbar "github.com/schollz/progressbar/v3"

	"github.com/hypha-dao/dao-contracts/dao-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
)

var newenvCmd = &cobra.Command{
	Use:   "newenv",
	Short: "creates a new DAO environment",
	Long:  "creates a new smart contract DAO environment",
	Run: func(cmd *cobra.Command, args []string) {

		env := setupEnvironment()
		fmt.Println(env.String())
		fmt.Println("\nDAO Environment Setup complete")
	},
}

func init() {
	RootCmd.AddCommand(newenvCmd)
}

var exchangeWasm, exchangeAbi string

// uses this hard-coded endpoint, does not try to use the config file
const testingEndpoint = "https://testnet.telos.caleos.io"
const creator = "hypha"

// const testingEndpoint = "http://localhost:8888"
// const creator = "eosio"

type member struct {
	Member eos.AccountName
	Doc    docgraph.Document
}

type environment struct {
	ctx context.Context
	api eos.API

	DAO           eos.AccountName
	HusdToken     eos.AccountName
	HyphaToken    eos.AccountName
	HvoiceToken   eos.AccountName
	SeedsToken    eos.AccountName
	Bank          eos.AccountName
	SeedsEscrow   eos.AccountName
	SeedsExchange eos.AccountName
	Events        eos.AccountName
	TelosDecide   eos.AccountName
	Whale         member
	Root          docgraph.Document

	VotingDurationSeconds int64
	HyphaDeferralFactor   int64
	SeedsDeferralFactor   int64

	NumPeriods     int
	PeriodDuration time.Duration

	PeriodPause        time.Duration
	VotingPause        time.Duration
	ChainResponsePause time.Duration

	Members []member
}

func envHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Variable"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}
}

func (e *environment) String() string {
	table := simpletable.New()
	table.Header = envHeader()

	kvs := make(map[string]string)
	kvs["DAO"] = string(e.DAO)
	kvs["HUSD Token"] = string(e.HusdToken)
	kvs["HVOICE Token"] = string(e.HvoiceToken)
	kvs["HYPHA Token"] = string(e.HyphaToken)
	kvs["SEEDS Token"] = string(e.SeedsToken)
	kvs["Bank"] = string(e.Bank)
	kvs["Escrow"] = string(e.SeedsEscrow)
	kvs["Exchange"] = string(e.SeedsExchange)
	kvs["Telos Decide"] = string(e.TelosDecide)
	kvs["Whale"] = string(e.Whale.Member)
	kvs["Voting Duration (s)"] = strconv.Itoa(int(e.VotingDurationSeconds))
	kvs["HYPHA deferral X"] = strconv.Itoa(int(e.HyphaDeferralFactor))
	kvs["SEEDS deferral X"] = strconv.Itoa(int(e.SeedsDeferralFactor))

	for key, value := range kvs {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: key},
			{Align: simpletable.AlignRight, Text: value},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}

func createAccountFromString(ctx context.Context, api *eos.API, accountName, privateKey string) (eos.AccountName, error) {

	key, err := ecc.NewPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("privateKey parameter is not a valid format: %s", err)
	}

	err = api.Signer.ImportPrivateKey(ctx, privateKey)
	if err != nil {
		return "", fmt.Errorf("Error importing key: %s", err)
	}

	return createAccount(ctx, api, accountName, key.PublicKey())
}

func addCodePerm(ctx context.Context, api *eos.API, accountName string, publicKey ecc.PublicKey) (eos.AccountName, error) {
	acct := toAccount(accountName, "account to create")

	codePermissionActions := []*eos.Action{system.NewUpdateAuth(acct,
		"active",
		"owner",
		eos.Authority{
			Threshold: 1,
			Keys: []eos.KeyWeight{{
				PublicKey: publicKey,
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
		}, "owner")}

	_, err := eostest.ExecTrx(ctx, api, codePermissionActions)
	if err != nil {
		return "", fmt.Errorf("Error filling tx opts: %s", err)
	}
	return acct, nil
}

func createAccount(ctx context.Context, api *eos.API, accountName string, publicKey ecc.PublicKey) (eos.AccountName, error) {
	acct := toAccount(accountName, "account to create")

	actions := []*eos.Action{system.NewNewAccount(creator, acct, publicKey)}
	_, err := eostest.ExecTrx(ctx, api, actions)
	if err != nil {
		return eos.AccountName(""), fmt.Errorf("Error filling tx opts: %s", err)
	}

	return addCodePerm(ctx, api, accountName, publicKey)
}

type tokenCreate struct {
	Issuer    eos.AccountName
	MaxSupply eos.Asset
}

func setupEnvironment() *environment {

	// t := testing.T{}

	daoHome := "/Users/max/dev/hypha/dao"
	// daoPrefix := daoHome + "/build/dao/dao."

	// artifactsHome := daoHome + "/dao-go/artifacts"
	// decidePrefix := artifactsHome + "/decide/decide."
	// treasuryPrefix := artifactsHome + "/treasury/treasury."
	// monitorPrefix := artifactsHome + "/monitor/monitor."
	// escrowPrefix := artifactsHome + "/escrow/escrow."

	exchangeWasm = daoHome + "/dao-go/mocks/seedsexchg/build/seedsexchg/seedsexchg.wasm"
	exchangeAbi = daoHome + "/dao-go/mocks/seedsexchg/build/seedsexchg/seedsexchg.abi"

	// Private key: 5KCZ9VBJMMiLaAY24Ro66mhx4vU1VcJELZVGrJbkUBATyqxyYmj
	// Public key: EOS5Kt3doXqfT6kVaAaRWWUWcANNXFUbgyuqjvpkhJtSJnBbkoAhQ
	creationPrivateKey := "5KCZ9VBJMMiLaAY24Ro66mhx4vU1VcJELZVGrJbkUBATyqxyYmj"
	// creationPublicKey := "EOS5Kt3doXqfT6kVaAaRWWUWcANNXFUbgyuqjvpkhJtSJnBbkoAhQ"
	// pubk, _ := ecc.NewPublicKey(creationPublicKey)

	var env environment
	env.DAO = eos.AN("dao1.hypha")
	env.HusdToken = eos.AN("token1.hypha")
	env.Bank = eos.AN("bank1.hypha")
	env.HyphaToken = eos.AN("token1.hypha")
	env.SeedsEscrow = eos.AN("escro1.hypha")
	env.SeedsExchange = eos.AN("tlosto.seeds")
	env.TelosDecide = eos.AN("td1.hypha")
	env.Events = eos.AN("publs1.hypha")

	env.api = *eos.New(testingEndpoint)
	// api.Debug = true
	env.ctx = context.Background()

	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(env.ctx, creationPrivateKey)
	err = keyBag.ImportPrivateKey(env.ctx, "5KBBECPRDtgPJNsTVUerFfxXpya1Ce93avTye9rn4oButkCSbPZ") //eostest.DefaultKey())
	if err != nil {
		panic(err)
	}
	env.api.SetSigner(keyBag)

	// _, err = addCodePerm(env.ctx, &env.api, "dao1.hypha", pubk)
	// _, err = addCodePerm(env.ctx, &env.api, "token1.hypha", pubk)
	// addCodePerm(env.ctx, &env.api, "bank1.hypha", pubk)
	// addCodePerm(env.ctx, &env.api, "escro1.hypha", pubk)
	// addCodePerm(env.ctx, &env.api, "td1.hypha", pubk)
	// addCodePerm(env.ctx, &env.api, "publs1.hypha", pubk)

	env.VotingDurationSeconds = 3600
	env.SeedsDeferralFactor = 100
	env.HyphaDeferralFactor = 25

	env.PeriodDuration, _ = time.ParseDuration("168h")
	env.NumPeriods = 100

	// bankPermissionActions := []*eos.Action{system.NewUpdateAuth(env.Bank,
	// 	"active",
	// 	"owner",
	// 	eos.Authority{
	// 		Threshold: 1,
	// 		Keys: []eos.KeyWeight{{
	// 			PublicKey: pubk,
	// 			Weight:    1,
	// 		}},
	// 		Accounts: []eos.PermissionLevelWeight{
	// 			{
	// 				Permission: eos.PermissionLevel{
	// 					Actor:      env.Bank,
	// 					Permission: "eosio.code",
	// 				},
	// 				Weight: 1,
	// 			},
	// 			{
	// 				Permission: eos.PermissionLevel{
	// 					Actor:      env.DAO,
	// 					Permission: "eosio.code",
	// 				},
	// 				Weight: 1,
	// 			}},
	// 		Waits: []eos.WaitWeight{},
	// 	}, "owner")}

	// _, err = eostest.ExecTrx(env.ctx, &env.api, bankPermissionActions)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Deploying DAO contract to 		: ", env.DAO)
	// _, err = eostest.SetContract(env.ctx, &env.api, env.DAO, daoPrefix+"wasm", daoPrefix+"abi")

	// fmt.Println("Deploying Treasury contract to 		: ", env.Bank)
	// _, err = eostest.SetContract(env.ctx, &env.api, env.Bank, treasuryPrefix+"wasm", treasuryPrefix+"abi")

	// fmt.Println("Deploying Escrow contract to 		: ", env.SeedsEscrow)
	// _, err = eostest.SetContract(env.ctx, &env.api, env.SeedsEscrow, escrowPrefix+"wasm", escrowPrefix+"abi")

	// fmt.Println("Deploying SeedsExchange contract to 		: ", env.SeedsExchange)
	// _, err = eostest.SetContract(env.ctx, &env.api, env.SeedsExchange, exchangeWasm, exchangeAbi)

	// fmt.Println("Deploying Events contract to 		: ", env.Events)
	// _, err = eostest.SetContract(env.ctx, &env.api, env.Events, monitorPrefix+"wasm", monitorPrefix+"abi")

	// // _, err = dao.CreateRoot(env.ctx, &env.api, env.DAO)

	// env.Root, err = docgraph.LoadDocument(env.ctx, &env.api, env.DAO, "cb7574fd1cfbdcbede5c8d7860ae5772a69785211cd56a52cb31e7ed492d60fb")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dao.NewTreasury(env.ctx, &env.api, env.TelosDecide, env.DAO)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Setting configuration options on DAO 		: ", env.DAO)
	// _, err = dao.SetIntSetting(env.ctx, &env.api, env.DAO, "voting_duration_sec", env.VotingDurationSeconds)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dao.SetIntSetting(env.ctx, &env.api, env.DAO, "seeds_deferral_factor_x100", env.SeedsDeferralFactor)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dao.SetIntSetting(env.ctx, &env.api, env.DAO, "hypha_deferral_factor_x100", env.HyphaDeferralFactor)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dao.SetIntSetting(env.ctx, &env.api, env.DAO, "paused", 0)
	// if err != nil {
	// 	panic(err)

	// }

	env.SeedsToken = eos.AN("token.seeds")

	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "hypha_token_contract", env.HyphaToken)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "hvoice_token_contract", env.HvoiceToken)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "husd_token_contract", env.HusdToken)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "seeds_token_contract", env.SeedsToken)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "seeds_escrow_contract", env.SeedsEscrow)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "publisher_contract", env.Events)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "treasury_contract", env.Bank)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "telos_decide_contract", env.TelosDecide)
	// dao.SetNameSetting(env.ctx, &env.api, env.DAO, "last_ballot_id", "hypha......1")

	// fmt.Println("Adding "+strconv.Itoa(env.NumPeriods)+" periods with duration 		: ", env.PeriodDuration)
	// _, err = dao.AddPeriods(env.ctx, &env.api, env.DAO, env.NumPeriods, env.PeriodDuration)
	// if err != nil {
	// 	panic(err)
	// }

	// deploy TD contract
	// fmt.Println("Deploying/configuring Telos Decide contract 		: ", env.TelosDecide)
	// _, err = eostest.SetContract(env.ctx, &env.api, env.TelosDecide, decidePrefix+"wasm", decidePrefix+"abi")

	// hvoiceMaxSupply, _ := eos.NewAssetFromString("1000000000.00 HVOICE")
	// _, err = dao.InitTD(env.ctx, &env.api, env.TelosDecide)
	// if err != nil {
	// 	panic(err)
	// }

	// transfer
	// _, err = dao.Transfer(env.ctx, &env.api, tlosToken, env.DAO, env.TelosDecide, tlosMaxSupply, "deposit")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dao.NewTreasury(env.ctx, &env.api, env.TelosDecide, env.DAO)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = dao.RegVoter(env.ctx, &env.api, env.TelosDecide, env.DAO)
	// if err != nil {
	// 	panic(err)
	// }

	// daoTokens, _ := eos.NewAssetFromString("10.00 HVOICE")
	// _, err = dao.Mint(env.ctx, &env.api, env.TelosDecide, env.DAO, env.DAO, daoTokens)
	// if err != nil {
	// 	panic(err)
	// }

	// // whaleTokens, _ := eos.NewAssetFromString("100.00 HVOICE")
	// // env.Whale, err = setupMember(env.ctx, &env.api, env.DAO, env.TelosDecide, "whale", whaleTokens)
	// // if err != nil {
	// // 	panic(err)
	// // }

	// // index := 1
	// // for index < 5 {

	// // 	memberNameIn := "dao1member" + strconv.Itoa(index)

	// setupMember(env.ctx, &env.api, env.DAO, env.TelosDecide, "mem1.hypha", daoTokens)
	// setupMember(env.ctx, &env.api, env.DAO, env.TelosDecide, "mem2.hypha", daoTokens)
	// setupMember(env.ctx, &env.api, env.DAO, env.TelosDecide, "mem3.hypha", daoTokens)
	// setupMember(env.ctx, &env.api, env.DAO, env.TelosDecide, "mem4.hypha", daoTokens)
	// setupMember(env.ctx, &env.api, env.DAO, env.TelosDecide, "mem5.hypha", daoTokens)

	// 	env.Members = append(env.Members, newMember)
	// 	index++
	// }

	return &env
}

func setupMember(ctx context.Context, api *eos.API,
	contract, telosDecide eos.AccountName, memberName string, hvoice eos.Asset) (member, error) {

	fmt.Println("Creating and enrolling new member  		: ", memberName, " 	with voting power	: ", hvoice.String())

	memberAccount := eos.AN(memberName)
	_, err := dao.RegVoter(ctx, api, telosDecide, memberAccount)
	if err != nil {
		panic(err)
	}

	_, err = dao.Mint(ctx, api, telosDecide, contract, memberAccount, hvoice)
	if err != nil {
		panic(err)
	}

	_, err = dao.Apply(ctx, api, contract, memberAccount, "apply to DAO")
	if err != nil {
		panic(err)
	}

	_, err = dao.Enroll(ctx, api, contract, contract, memberAccount)
	if err != nil {
		panic(err)
	}

	pause(time.Second, "Build block...", "")

	memberDoc, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, "member")
	if err != nil {
		panic(err)
	}

	return member{
		Member: memberAccount,
		Doc:    memberDoc,
	}, nil
}

func pause(seconds time.Duration, headline, prefix string) {
	if headline != "" {
		fmt.Println(headline)
	}

	bar := progressbar.NewOptions(100,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(90),
		// progressbar.OptionShowIts(),
		progressbar.OptionSetDescription("[cyan]"+fmt.Sprintf("%20v", prefix)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	chunk := seconds / 100
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(chunk)
	}
	fmt.Println()
	fmt.Println()
}
