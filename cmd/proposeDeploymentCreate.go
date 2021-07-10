package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const proposalNamePromptLabel = "Required: name of the proposal (eosio::name format)"
const commitPromptLabel = "Required: commit hash of the github repo to use for deployment"
const notesPromptLabel = "Optional: notes to attach to the deployment proposal document"
const developerPromptLabel = "Optional: account name of the developer who contributed most to the upgrade"
const existingDocumentLabel = "Optional: hash of another document within the DAO that describes the deployment, e.g. an approved policy doc"
const accountPromptLabel = "Required: name of account to deploy contract to"

func grabInput(field, promptLabel string) (string, error) {
	if len(viper.GetString(field)) == 0 {
		fmt.Println()
		prompt := promptui.Prompt{Label: promptLabel}
		result, err := prompt.Run()
		if err != nil {
			return string(""), fmt.Errorf("cannot capture input: %v %v", field, err)
		}
		return result, nil
	} else {
		return viper.GetString(field), nil
	}
}

var proposeDeploymentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "proposes a contract deployment based on a git commit",
	Long:  "proposes a contract deployment based on a git commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		proposalName, err := grabInput("propose-deployment-create-cmd-proposal-name", proposalNamePromptLabel)
		if err != nil {
			return fmt.Errorf("cannot get input: %v %v", proposalNamePromptLabel, err)
		}

		account, err := grabInput("propose-deployment-create-cmd-account", accountPromptLabel)
		if err != nil {
			return fmt.Errorf("cannot get input: %v %v", accountPromptLabel, err)
		}

		// TODO: query list of recent commits from github
		commit, err := grabInput("propose-deployment-create-cmd-commit", commitPromptLabel)
		if err != nil {
			return fmt.Errorf("cannot get input: %v %v", commitPromptLabel, err)
		}

		developer, err := grabInput("propose-deployment-create-cmd-developer", developerPromptLabel)
		if err != nil {
			return fmt.Errorf("cannot get input: %v %v", developerPromptLabel, err)
		}

		document, err := grabInput("propose-deployment-create-cmd-document", existingDocumentLabel)
		if err != nil {
			return fmt.Errorf("cannot get input: %v %v", existingDocumentLabel, err)
		}

		notes, err := grabInput("propose-deployment-create-cmd-notes", notesPromptLabel)
		if err != nil {
			return fmt.Errorf("cannot get input: %v %v", notesPromptLabel, err)
		}

		accountToDeploy := eos.AN(account)

		dir, err := os.MkdirTemp(os.TempDir(), "daoctl_propose_deployment")
		if err != nil {
			return fmt.Errorf("cannot create a temporary directory: %v", err)
		}
		zlog.Info("creating temp directory to build contract: " + dir)

		repo, err := git.PlainClone(dir, false, &git.CloneOptions{
			URL:               viper.GetString("DAORepo"),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Progress:          os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("cannot clone repo: %v", err)
		}
		buildDir := dir + "/build"

		err = os.Mkdir(buildDir, 0700)
		if err != nil {
			return fmt.Errorf("cannot create build directory: %v", err)
		}

		ref, err := repo.Head()
		if err != nil {
			return fmt.Errorf("cannot access the repo HEAD: %v", err)
		}
		zlog.Info("current repo HEAD " + ref.Hash().String())

		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("cannot access the repo Worktree: %v", err)
		}
		zlog.Info("executing a checkout of commit: " + commit)

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commit),
		})
		if err != nil {
			return fmt.Errorf("unable to checkout to commit %v", err)
		}

		cmake := exec.Command("cmake", dir)
		cmake.Dir = buildDir
		zlog.Info("running cmake - " + cmake.String())
		cmake.Run()

		make := exec.Command("make", "-j"+strconv.Itoa(runtime.NumCPU()))
		make.Dir = buildDir
		zlog.Info("running make to build contracts - " + make.String())
		make.Run()

		hProp := hyphaProposal{}
		eProp := eosioProposal{}

		hProp.ProposalName = eos.Name(proposalName)
		hProp.Proposer = eos.AccountName(viper.GetString("DAOUser"))

		eProp.ProposalName = eos.Name(proposalName)
		eProp.Proposer = eos.AccountName(viper.GetString("DAOUser"))
		eProp.RequestedApprovals = []eos.PermissionLevel{
			{
				Actor:      "gh.hypha",
				Permission: "active",
			},
			{
				Actor:      "m.hypha",
				Permission: "active",
			},
			{
				Actor:      "j.hypha",
				Permission: "active",
			},
			{
				Actor:      "jj.hypha",
				Permission: "active",
			},
			{
				Actor:      "l.hypha",
				Permission: "active",
			},
		}

		hProp.ContentGroups = []docgraph.ContentGroup{{
			docgraph.ContentItem{
				Label: "content_group_name",
				Value: &docgraph.FlexValue{
					BaseVariant: eos.BaseVariant{
						TypeID: docgraph.GetVariants().TypeID("string"),
						Impl:   "Deployment Proposal Details",
					},
				},
			},
			docgraph.ContentItem{
				Label: "notes",
				Value: &docgraph.FlexValue{
					BaseVariant: eos.BaseVariant{
						TypeID: docgraph.GetVariants().TypeID("string"),
						Impl:   notes,
					},
				},
			},
			docgraph.ContentItem{
				Label: "document",
				Value: &docgraph.FlexValue{
					BaseVariant: eos.BaseVariant{
						TypeID: docgraph.GetVariants().TypeID("string"), // TODO: check if valid hash?
						Impl:   document,
					},
				},
			},
			docgraph.ContentItem{
				Label: "github_commit",
				Value: &docgraph.FlexValue{
					BaseVariant: eos.BaseVariant{
						TypeID: docgraph.GetVariants().TypeID("string"),
						Impl:   "https://github.com/hypha-dao/dao-contracts/commit/" + commit,
					},
				},
			},
			docgraph.ContentItem{
				Label: "developer",
				Value: &docgraph.FlexValue{
					BaseVariant: eos.BaseVariant{
						TypeID: docgraph.GetVariants().TypeID("name"),
						Impl:   eos.Name(developer),
					},
				},
			},
		}}

		setCodeAction, err := system.NewSetCode(accountToDeploy, buildDir+"/dao/dao.wasm")
		if err != nil {
			return fmt.Errorf("unable construct set_code action: %v", err)
		}

		setAbiAction, err := system.NewSetABI(accountToDeploy, buildDir+"/dao/dao.abi")
		if err != nil {
			return fmt.Errorf("unable construct set_abi action: %v", err)
		}

		txOpts := &eos.TxOptions{}
		if err := txOpts.FillFromChain(ctx, api); err != nil {
			return fmt.Errorf("error filling tx opts: %s", err)
		}

		setCodeActions := []*eos.Action{setCodeAction, setAbiAction}
		eProp.Transaction = eos.NewTransaction(setCodeActions, txOpts)

		oneHour, _ := time.ParseDuration("1h")
		eProp.Transaction.SetExpiration(oneHour * 24 * 7)

		actions := []*eos.Action{
			{
				Account: eos.AN("eosio.msig"),
				Name:    eos.ActN("propose"),
				Authorization: []eos.PermissionLevel{
					{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
				},
				ActionData: eos.NewActionData(eProp),
			},
			{
				Account: eos.AN(viper.GetString("MsigContract")),
				Name:    eos.ActN("propose"),
				Authorization: []eos.PermissionLevel{
					{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
				},
				ActionData: eos.NewActionData(hProp),
			},
		}

		// msigTrx, err := json.MarshalIndent(d, "", "  ")
		// if err != nil {
		// 	return fmt.Errorf("cannot marshal object to json: %s", err)
		// }

		// _ = ioutil.WriteFile("msig-transaction.json", msigTrx, 0644)

		pushEOSCActions(ctx, api, actions[0], actions[1])
		return nil
	},
}

func init() {
	proposeDeploymentCmd.AddCommand(proposeDeploymentCreateCmd)
	proposeDeploymentCreateCmd.Flags().StringP("proposal-name", "", "", proposalNamePromptLabel)
	proposeDeploymentCreateCmd.Flags().StringP("commit", "", "", commitPromptLabel)
	proposeDeploymentCreateCmd.Flags().StringP("notes", "n", "", notesPromptLabel)
	proposeDeploymentCreateCmd.Flags().StringP("developer", "d", "", developerPromptLabel)
	proposeDeploymentCreateCmd.Flags().StringP("document", "", "", existingDocumentLabel)
	proposeDeploymentCreateCmd.Flags().StringP("account", "", "", accountPromptLabel)
}
