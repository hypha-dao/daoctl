package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCalendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "print the calendar",
	Long:  "print a table with each of the time periods",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		// gc, err := util.GetCache(ctx, api, contract)
		// if err != nil {
		// 	return fmt.Errorf("cannot get cache: %v", err)
		// }

		// rootDocument, err := docgraph.LoadDocument(ctx, api, contract, viper.GetString("RootNode"))
		// if err != nil {
		// 	return fmt.Errorf("cannot load root document: %v", err)
		// }
		// fmt.Println(rootDocument)

		// startEdges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, rootDocument, eos.Name("start"))
		// if err != nil {
		// 	return fmt.Errorf("error while retrieving start edge: %v", err)
		// }
		// fmt.Println(startEdges)
		// if len(startEdges) == 0 {
		// 	return fmt.Errorf("no start edge from the root node exists: %v", err)
		// }

		startPeriod := viper.GetString("CalendarStart") //"7706e72c29af438f309a99391fa8e8e3dcef0db438d0d24daf6fc4cf29697bff"

		startPeriodDoc, err := docgraph.LoadDocument(ctx, api, contract, startPeriod) //startEdges[0].ToNode.String())

		// startPeriodDoc, err := docgraph.LoadDocument(ctx, api, contract, startEdges[0].ToNode.String())
		if err != nil {
			return fmt.Errorf("error loading the start period document: %v", err)
		}

		period, err := models.NewPeriod(ctx, api, contract, startPeriodDoc)
		if err != nil {
			return fmt.Errorf("cannot convert document to period type: %v", err)
		}

		periodTable := views.PeriodTable(period)
		periodTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Println("\n" + periodTable.String() + "\n\n")

		return nil
	},
}

func init() {
	getCmd.AddCommand(getCalendarCmd)
}
