package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type deployment struct {
	Proposer           eos.AccountName         `json:"proposer"`
	ProposalName       eos.Name                `json:"proposal_name"`
	RequestedApprovals []eos.PermissionLevel   `json:"requested"`
	ContentGroups      []docgraph.ContentGroup `json:"content_groups"`
	Transaction        *eos.Transaction        `json:"trx"`
}

var proposeDeploymentCmd = &cobra.Command{
	Use:   "deployment [commit]",
	Short: "proposes a contract deployment based on a git commit",
	Long:  "proposes a contract deployment based on a git commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		commit := "98fc0294c7e415a468afb975d4d7948b46179f43" //args[2]
		accountToDeploy := eos.AN("dao1.hypha")
		notes := "these notes describe the deployment"
		developer := eos.AN("m.hypha")

		dir, err := os.MkdirTemp(os.TempDir(), "daoctl_propose_deployment")
		if err != nil {
			return fmt.Errorf("cannot create a temporary directory: %v", err)
		}
		zap.S().Info("creating temp directory to build contract: " + dir)

		repo, err := git.PlainClone(dir, false, &git.CloneOptions{
			URL:               "https://github.com/hypha-dao/dao-contracts",
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
			return fmt.Errorf("cannot create a temporary directory: %v", err)
		}
		zap.S().Info("printing the repo HEAD " + ref.Hash().String())

		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("cannot create a temporary directory: %v", err)
		}

		zap.S().Info("git checkout " + commit)

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commit),
		})
		if err != nil {
			return fmt.Errorf("unable to checkout to commit %v", err)
		}

		sub, err := w.Submodule("document-graph")
		if err != nil {
			return fmt.Errorf("cannot get document-graph submodule repo: %v %v", buildDir, err)
		}

		sr, err := sub.Repository()
		if err != nil {
			return fmt.Errorf("cannot get document-graph repo: %v %v", buildDir, err)
		}

		sw, err := sr.Worktree()
		if err != nil {
			return fmt.Errorf("cannot get document-graph sw: %v %v", buildDir, err)
		}

		zap.S().Info("running document-graph submodule update --remote")
		err = sw.Pull(&git.PullOptions{
			RemoteName: "origin",
		})
		if err != nil {
			return fmt.Errorf("cannot get document-graph sw: %v %v", buildDir, err)
		}

		cmake := exec.Command("cmake", dir)
		cmake.Dir = buildDir
		zap.S().Info("running cmake - " + cmake.String())
		cmake.Run()

		make := exec.Command("make", "-j"+strconv.Itoa(runtime.NumCPU()))
		make.Dir = buildDir
		zap.S().Info("running make to build contracts - " + make.String())
		make.Run()

		files, err := ioutil.ReadDir(buildDir)
		if err != nil {
			return fmt.Errorf("cannot read build directory: %v %v", buildDir, err)
		}

		zap.S().Info("listing files from build directory : " + buildDir)
		for _, file := range files {
			fmt.Println(file.Name())
		}

		ref, err = repo.Head()
		if err != nil {
			return fmt.Errorf("unable to switch to Head: %v", err)
		}

		zap.S().Info("printing the repo HEAD " + ref.Hash().String())
		d := deployment{}
		d.ProposalName = "deployment"
		d.Proposer = eos.AccountName(viper.GetString("DAOUser"))
		d.RequestedApprovals = []eos.PermissionLevel{
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
		d.ContentGroups = []docgraph.ContentGroup{{
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
						Impl:   developer,
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
		d.Transaction = eos.NewTransaction(setCodeActions, txOpts)

		actions := []*eos.Action{{
			Account: eos.AN("msig.hypha"),
			Name:    eos.ActN("propose"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(d),
		}}

		// msigTrx, err := json.MarshalIndent(d, "", "  ")
		// if err != nil {
		// 	return fmt.Errorf("cannot marshal object to json: %s", err)
		// }

		// _ = ioutil.WriteFile("msig-transaction.json", msigTrx, 0644)

		pushEOSCActions(ctx, api, actions[0])
		return nil
	},
}

func init() {
	proposeCmd.AddCommand(proposeDeploymentCmd)
}
