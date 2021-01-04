package cmd

import (
	"context"
	"fmt"
	"image"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/manifoldco/promptui"
	"github.com/patrickmn/go-cache"
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
		fmt.Sprintf("Created Date|%v", p.Primary.CreatedDate.Time.Format("2006 Jan 02 15:04:05")),
	}

	fmt.Println(columnize.SimpleFormat(output))
	fmt.Println()
	printContentGroups(p)
	// printEdges(ctx, api, p)
}

// Page ...
type Page struct {
	Primary     docgraph.Document
	FromEdges   []docgraph.Edge
	ToEdges     []docgraph.Edge
	EdgePrompts []edgeChoice
	Graph       image.Image
}

func (p *Page) getKey() string {
	return p.Primary.Hash.String()
}

func getLabel(d *docgraph.Document) string {
	documentLabel, _ := d.GetContent("node_label")
	if documentLabel != nil {
		if len(documentLabel.String()) <= 40 {
			return documentLabel.String()
		}
		return documentLabel.String()[:40]
	}
	return "unlabeled"
}

func getType(d *docgraph.Document) string {
	documentType, _ := d.GetContent("type")
	if documentType != nil {
		return documentType.String()
	}
	return "untyped"
}

type edgeChoice struct {
	Name        eos.Name
	Forward     bool
	To          string
	ToType      string
	ToLabel     string
	CreatedDate eos.BlockTimestamp
	Column1     string // First 69 characters
	Column2     string //
}

func (e *edgeChoice) GetWidth() int {
	return 25
}

func getPage(ctx context.Context, api *eos.API, pageCache, documentCache *cache.Cache, contract eos.AccountName, hash string) Page {
	pager, found := pageCache.Get(hash)
	if found {
		return pager.(Page)
	}

	var err error
	page := Page{}
	page.Primary = getDocument(ctx, api, documentCache, contract, hash)

	page.FromEdges, err = docgraph.GetEdgesFromDocument(ctx, api, contract, page.Primary)
	if err != nil {
		log.Println("ERROR: Cannot get edges from document: ", err)
	}

	page.ToEdges, err = docgraph.GetEdgesToDocument(ctx, api, contract, page.Primary)
	if err != nil {
		log.Println("ERROR: Cannot get edges to document: ", err)
	}

	page.EdgePrompts = make([]edgeChoice, len(page.FromEdges)+len(page.ToEdges))

	for i, edge := range page.FromEdges {

		document := getDocument(ctx, api, documentCache, contract, edge.ToNode.String())

		page.EdgePrompts[i] = edgeChoice{
			Name:        edge.EdgeName,
			Forward:     true,
			To:          edge.ToNode.String(),
			ToType:      getType(&document),
			ToLabel:     getLabel(&document),
			CreatedDate: edge.CreatedDate,
			Column1:     string("                                                                          [x]"),
			Column2: string(" ---> " + fmt.Sprintf("%-12v", string(edge.EdgeName)) +
				" ---> " + getLabel(&document) + " (" + getType(&document) + ")"),
		}
	}

	for i, edge := range page.ToEdges {

		document := getDocument(ctx, api, documentCache, contract, edge.FromNode.String())

		page.EdgePrompts[i+len(page.FromEdges)] = edgeChoice{
			Name:        edge.EdgeName,
			Forward:     false,
			To:          edge.FromNode.String(),
			ToType:      getType(&document),
			ToLabel:     getLabel(&document),
			CreatedDate: edge.CreatedDate,
			Column1: string(
				fmt.Sprintf("%49v", getLabel(&document)+
					" ("+
					getType(&document)) +
					") ---> " +
					fmt.Sprintf("%12v", string(edge.EdgeName)) +
					" ---> [x]"),
			Column2: string("                                                                          "),
		}
	}

	sort.SliceStable(page.EdgePrompts, func(i, j int) bool {
		return page.EdgePrompts[i].CreatedDate.Before(page.EdgePrompts[j].CreatedDate.Time)
	})

	pageCache.Set(page.Primary.Hash.String(), page, cache.DefaultExpiration)

	return page
}

func getDocument(ctx context.Context, api *eos.API, c *cache.Cache, contract eos.AccountName, hash string) docgraph.Document {

	documenter, found := c.Get(hash)
	if found {
		return documenter.(docgraph.Document)
	}

	document, err := docgraph.LoadDocument(ctx, api, contract, hash)
	if err != nil {
		log.Println("Document not found: " + hash)
		return docgraph.Document{}
	}
	c.Set(document.Hash.String(), document, cache.DefaultExpiration)

	return document
}

func loadCache(ctx context.Context, api *eos.API, pages, documents *cache.Cache, contract eos.AccountName, startingNode string) {

	go func() {
		page := getPage(ctx, api, pages, documents, contract, startingNode)

		for _, edge := range page.ToEdges {
			getPage(ctx, api, pages, documents, contract, edge.ToNode.String())
		}

		for _, edge := range page.FromEdges {
			getPage(ctx, api, pages, documents, contract, edge.ToNode.String())
		}
	}()
}

var getDocumentCmd = &cobra.Command{
	Use:   "document [hash]",
	Short: "retrieve document details and navigate the graph",
	Long:  "retrieve the detailed content within a document",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		var hash string

		if viper.GetBool("get-document-cmd-last") {
			lastDocument, err := docgraph.GetLastDocument(ctx, api, contract)
			if err != nil {
				panic(fmt.Errorf("cannot get last document: %v", err))
			}
			hash = lastDocument.Hash.String()
			fmt.Println("Last document hash: " + hash)
		}

		if len(args) == 1 {
			hash = args[0]
		}

		fmt.Println("STARTING hash: " + hash)

		var page Page
		pages := cache.New(5*time.Minute, 10*time.Minute)
		documents := cache.New(5*time.Minute, 10*time.Minute)

		for {

			page = getPage(ctx, api, pages, documents, contract, hash)

			loadCache(ctx, api, pages, documents, contract, hash)

			printDocument(ctx, api, &page)

			fmt.Println("                          ")
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}  {{ \"from: node_label (type) --->    edge_name ---> [X] ---> edge_name    ---> to: node_label (type) \" | faint }}",
				Active:   "{{ .Column1 | cyan}}{{ .Column2 | cyan}}",
				Inactive: "{{ .Column1 }}{{ .Column2 }}",
				Selected: "{{ \"  -- selected -- \" | yellow}}{{ .Column1 | yellow}}{{ .Column2 | yellow}} {{ \"  -- selected -- \" | yellow}}",
			}

			searcher := func(input string, index int) bool {
				pepper := page.EdgePrompts[index]
				name := strings.Replace(strings.ToLower(string(pepper.Name)), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			// prompt := promptui.Prompt{
			// 	Label:     "Save Image",
			// 	IsConfirm: true,
			// }

			// saveImage, err := prompt.Run()

			// fmt.Println("You selected: " + saveImage)
			// if err != nil {
			// 	fmt.Printf("Prompt failed %v\n", err)
			// 	return
			// }

			// if saveImage == "y" {
			// 	f, err := os.Create(page.Primary.Hash.String() + ".png")
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	defer f.Close()

			// 	err = png.Encode(f, page.Graph)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }

			prompt2 := promptui.Select{
				Label:     "Select an edge to navigate:",
				Items:     page.EdgePrompts,
				Templates: templates,
				Size:      len(page.FromEdges) + len(page.ToEdges),
				Searcher:  searcher,
			}

			i, _, err := prompt2.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			fmt.Println()
			fmt.Println("-------------------------------------------------------------------------------------------------")
			fmt.Println()

			hash = page.EdgePrompts[i].To

			// err = ioutil.WriteFile("last-doc.tmp", []byte(hash), 0644)
			// if err != nil {
			// 	fmt.Printf("Failed to write temporary file %v\n", err)
			// 	return
			// }
		}
	},
}

func init() {
	getDocumentCmd.Flags().BoolP("last", "l", false, "retrieve the most recently created document")
	getCmd.AddCommand(getDocumentCmd)
}

// saves a file PNG of the image, but commenting out to avoid the dependency for now
func saveImage(ctx context.Context, api *eos.API, pageCache, documentCache *cache.Cache, contract eos.AccountName, hash string) {

	// page := getPage(ctx, api, pageCache, documentCache, contract, hash)

	// var g *graphviz.Graphviz
	// g = graphviz.New()
	// graph, err := g.Graph()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() {
	// 	if err := graph.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	g.Close()
	// }()

	// primaryNode, err := graph.CreateNode(getLabel(&page.Primary))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for i, edge := range page.FromEdges {

	// 	document := getDocument(ctx, api, documentCache, contract, edge.ToNode.String())
	// 	secondaryNode, err := graph.CreateNode(getLabel(&document))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	spEdge, err := graph.CreateEdge(string(edge.EdgeName), secondaryNode, primaryNode)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	spEdge.SetLabel(string(edge.EdgeName))
	// }

	// for i, edge := range page.ToEdges {

	// 	document := getDocument(ctx, api, documentCache, contract, edge.FromNode.String())
	// 	secondaryNode, err := graph.CreateNode(getLabel(&document))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	psEdge, err := graph.CreateEdge(string(edge.EdgeName), primaryNode, secondaryNode)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	psEdge.SetLabel(string(edge.EdgeName))
	// }

	// page.Graph, err = g.RenderImage(graph)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
