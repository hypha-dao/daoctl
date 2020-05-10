package main

import (
  "github.com/hypha-dao/daoctl/cmd"
)

func main() {

	cmd.Execute()


}

// 	api := eos.New("https://api.telos.kitchen")
// 	// includeProposals := true
// 	// api := eos.New("https://")

// 	// infoResp, _ := api.GetInfo(ctx)
// 	// infoRespStr, _ := json.MarshalIndent(infoResp, "", "  ")
// 	periods := LoadPeriods(api)

// 	ctx := context.Background()
// 	roles := Roles(ctx, api, periods)
// 	// fmt.Println("\n\n" + RoleTable(roles) + "\n\n")

// 	// if includeProposals {
// 	// propRoles := ProposedRoles(ctx, api, periods)
// 	// fmt.Println("\n\n" + RoleTable(propRoles) + "\n\n")
// 	// }

// 	assignments := Assignments(ctx, api, roles, periods)
// 	assignmentsTable := AssignmentTable(assignments)
// 	assignmentsTable.SetStyle(simpletable.StyleCompactLite)
// 	fmt.Println("\n\n" + assignmentsTable.String() + "\n\n")

// 	data := TableToData(assignmentsTable)
// 	file, err := os.Create("assignments.csv")
// 	checkError("Cannot create file", err)
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	for _, value := range data {
// 		err := writer.Write(value)
// 		checkError("Cannot write to file", err)
// 	}
// }

// func checkError(message string, err error) {
// 	if err != nil {
// 		log.Fatal(message, err)
// 	}
// }



// PrintPayouts(context.Background(), api, periods, true)
