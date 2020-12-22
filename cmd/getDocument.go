package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/manifoldco/promptui"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cleanString(input string) string {
	input = strings.Replace(input, "\n", "", -1)

	if len(input) > 65 {
		return input[:40]
	}
	return input
}

func printContentGroups(p *Page) {

	fmt.Println("ContentGroups")
	for _, contentGroup := range p.Primary.ContentGroups {
		fmt.Println("  ContentGroup")

		for _, content := range contentGroup {
			fmt.Print("    ")
			fmt.Printf("%-35v", cleanString(content.Label))
			fmt.Printf("%-65v\n", cleanString(content.Value.String()))
		}
	}
	fmt.Println()
}

func printEdges(ctx context.Context, api *eos.API, p *Page) {

	colorCyan := "\033[36m"
	colorReset := "\033[0m"
	colorRed := "\033[31m"

	if len(p.ToEdges) > 0 {

		sort.SliceStable(p.ToEdges, func(i, j int) bool {
			return p.ToEdges[i].CreatedDate.Before(p.ToEdges[j].CreatedDate.Time)
		})

		toEdgesTable := views.EdgeTable(p.ToEdges, true, false)
		toEdgesTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Print(string(colorCyan), "\n                                                                                           to this node  ---------->")
		fmt.Println(string(colorReset))
		fmt.Println("\n" + toEdgesTable.String() + "\n\n")
	}

	if len(p.FromEdges) > 0 {

		sort.SliceStable(p.FromEdges, func(i, j int) bool {
			return p.FromEdges[i].CreatedDate.Before(p.FromEdges[j].CreatedDate.Time)
		})

		fromEdgesTable := views.EdgeTable(p.FromEdges, true, true)
		fromEdgesTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Print(string(colorCyan), "\n                                    from this node  ---------->")
		fmt.Println(string(colorReset))
		fmt.Println("\n" + fromEdgesTable.String() + "\n\n")
		// fmt.Print(string(colorCyan), "\n  <---   use 'rev <edge-name>' for reverse                                        use 'fwd <edge-name>' for forward --->")
		// fmt.Print(string(colorCyan), "\n  <edge-name> defaults to first in list if left blank")
		fmt.Println(string(colorReset))

	} else {
		fmt.Print(string(colorRed), "\n                                                                                                        ----------")
		fmt.Print(string(colorRed), "\n                                                                                                        |  END   |")
		fmt.Print(string(colorRed), "\n                                                                                                        ----------")
		fmt.Println(string(colorReset))
	}
}

func printDocument(ctx context.Context, api *eos.API, p *Page) {
	fmt.Println("Document Details")

	fmt.Println()
	output := []string{
		fmt.Sprintf("ID|%v", strconv.Itoa(int(p.Primary.ID))),
		fmt.Sprintf("Hash|%v", p.Primary.Hash.String()),
		fmt.Sprintf("Creator|%v", string(p.Primary.Creator)),
		fmt.Sprintf("Created Date|%v", p.Primary.CreatedDate.Time.Format("2006 Jan 02")),
	}

	fmt.Println(columnize.SimpleFormat(output))
	fmt.Println()
	printContentGroups(p)
	printEdges(ctx, api, p)
}

type edgeChoice struct {
	Name        eos.Name
	Forward     bool
	To          string
	ToType      string
	ToLabel     string
	CreatedDate eos.BlockTimestamp
}

// Page ...
type Page struct {
	Primary   docgraph.Document
	FromEdges []docgraph.Edge
	ToEdges   []docgraph.Edge
	Choices   []edgeChoice
}

var getDocumentCmd = &cobra.Command{
	Use:   "document [hash]",
	Short: "retrieve document details and navigate the graph",
	Long:  "retrieve the detailed content within a document",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		// var pages []Page

		var hash string
		if len(args) == 0 {
			hash = viper.GetString("RootNode")
		} else {
			hash = args[0]
		}

		for {

			page := Page{}
			// pages = append(pages)

			tempDocument, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), hash)
			if err != nil {
				panic("Document not found: " + hash)
			}
			page.Primary = tempDocument

			page.FromEdges, err = docgraph.GetEdgesFromDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), page.Primary)
			if err != nil {
				fmt.Println("ERROR: Cannot get edges from document: ", err)
			}

			page.ToEdges, err = docgraph.GetEdgesToDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), page.Primary)
			if err != nil {
				fmt.Println("ERROR: Cannot get edges to document: ", err)
			}

			printDocument(ctx, api, &page)

			edgePrompts := make([]edgeChoice, len(page.FromEdges)+len(page.ToEdges))
			for i, edge := range page.FromEdges {

				document, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), edge.ToNode.String())
				if err != nil {
					panic("Document not found: " + edge.ToNode.String())
				}

				typeLabel := "Unknown"
				documentType, _ := document.GetContent("type")
				if documentType != nil {
					// if isSkipped(documentType.String()) {
					// 	continue
					// }
					typeLabel = documentType.String()
				}

				nodeLabel := "Unknown"
				documentLabel, _ := document.GetContent("node_label")
				if documentLabel != nil {
					nodeLabel = documentLabel.String()
				}

				edgePrompts[i] = edgeChoice{
					Name:        edge.EdgeName,
					Forward:     true,
					To:          edge.ToNode.String(),
					ToType:      typeLabel,
					ToLabel:     nodeLabel,
					CreatedDate: edge.CreatedDate,
				}
			}

			for i, edge := range page.ToEdges {

				document, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), edge.FromNode.String())
				if err != nil {
					panic("Document not found: " + edge.ToNode.String())
				}

				typeLabel := "Unknown"
				documentType, _ := document.GetContent("type")
				if documentType != nil {
					// if isSkipped(documentType.String()) {
					// 	continue
					// }
					typeLabel = documentType.String()
				}

				nodeLabel := "Unknown"
				documentLabel, _ := document.GetContent("node_label")
				if documentLabel != nil {
					nodeLabel = documentLabel.String()
				}

				edgePrompts[i+len(page.FromEdges)] = edgeChoice{
					Name:        edge.EdgeName,
					Forward:     false,
					To:          edge.FromNode.String(),
					ToType:      typeLabel,
					ToLabel:     nodeLabel,
					CreatedDate: edge.CreatedDate,
				}
			}

			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U0001F9ED  {{if .Forward}}{{ .Name | cyan }} ---> {{ .ToLabel }} {{else}}{{ .ToLabel }} <--- {{ .Name | cyan }}{{end}}",
				Inactive: "    {{if .Forward}}{{ .Name | faint }} ---> {{ .ToLabel | faint }} {{else}}{{ .ToLabel | faint}} <--- {{ .Name | faint }}{{end}}",
				Selected: "\U0001F9ED  {{if .Forward}}{{ .Name | cyan }} ---> {{ .ToLabel }} {{else}}{{ .ToLabel }} <--- {{ .Name | cyan }}{{end}}",
				Details: `
		--------- Edge Details ----------
		{{ "Edge Name:" | faint }}	{{ .Name }}
		{{ "Node Label:" | faint }}	{{ .ToLabel }}
		{{ "Type:" | faint }}	{{ .ToType }}
		{{ "Forward:" | faint }}	{{ .Forward }}
		{{ "Created Date:" | faint }}	{{ .CreatedDate }}
		{{ "To:" | faint }}	{{ .To }}`,
			}

			searcher := func(input string, index int) bool {
				pepper := edgePrompts[index]
				name := strings.Replace(strings.ToLower(string(pepper.Name)), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			prompt2 := promptui.Select{
				Label:     "Select an edge",
				Items:     edgePrompts,
				Templates: templates,
				Size:      4,
				Searcher:  searcher,
			}

			i, _, err := prompt2.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			hash = edgePrompts[i].To
			fmt.Printf("You choose number %d: %s\n", i+1, string(edgePrompts[i].Name))

			err = ioutil.WriteFile("last-doc.tmp", []byte(hash), 0644)
			if err != nil {
				fmt.Printf("Failed to write temporary file %v\n", err)
				return
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getDocumentCmd)
}
