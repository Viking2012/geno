/*
Copyright © 2022 Alexander Orban <alexander.orban@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Viking2012/geno/geno"
	"github.com/Viking2012/geno/pkg"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	fPath              string
	refreshConstraints bool
	query              geno.Query
	constraints        geno.Constraints
	nodesFoundCount    map[string]int = make(map[string]int)
	relsFoundCount     map[string]int = make(map[string]int)
	nodesMergedCount   map[string]int = make(map[string]int)
	relsMergedCount    map[string]int = make(map[string]int)
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "import a json file of nodes and/or relationships",
	Long: `Import a json file containing nodes and/or relationships.
The file must be in the format:
{
	"nodes":[
		{"identity":1,"labels":["Label1","Label2",...],"Properties":{"Prop1":Value1,"Prop2":Value2,...}},
		{"identity":2,"labels":["Label3","Label4",...],"Properties":{"Prop1":Value3,"Prop2":Value4,...}},
	],
	"rels":[
		{"identity":1, "start":1,"end":2,"type":"Rel_Type","Properties":{"RelProp1":Value1, "RelProp2":Value2,...}}
	]
}`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if fPath == "" {
			return errors.New("filepath cannot be empty")
		}
		raw, err := os.ReadFile(fPath)
		if err != nil {
			return err
		}
		graph, err := pkg.GetGraphFromJson(raw)
		if err != nil {
			return err
		}
		driver, err := geno.NewDriver("neo4j://"+cfg.Server, neo4j.BasicAuth(cfg.User, cfg.GetPassword(), ""))
		if err != nil {
			return err
		}

		// fmt.Println(cfg.Constraints)
		if refreshConstraints {
			constraints, err = driver.GetConstraints(cfg.Database)
			if err != nil {
				return err
			}
		} else {
			constraints = cfg.Constraints[cfg.Database]
		}

		query = geno.NewQuery(&driver, &constraints)

		bar := progressbar.Default(int64(len(graph.Nodes)), "nodes")
		for _, node := range graph.Nodes {
			summary, err := query.MergeNode(cfg.Database, node)
			if err != nil {
				return err
			}
			bar.Add(1)
			for _, l := range node.Labels {
				nodesFoundCount[l]++
				nodesMergedCount[l] += summary.Counters().NodesCreated()
			}

		}
		bar = progressbar.Default(int64(len(graph.Relationships)), "rels ")
		for _, rel := range graph.Relationships {
			summary, err := query.MergeRelationship(cfg.Database, rel)
			if err != nil {
				return err
			}
			bar.Add(1)
			relsFoundCount[rel.Label]++
			relsMergedCount[rel.Label] += summary.Counters().RelationshipsCreated()
		}

		fmt.Println("nodes report:")
		for lab, cnt := range nodesFoundCount {
			fmt.Println("\tNode type:", lab, " found:", cnt, " merged:", nodesMergedCount[lab])
		}
		fmt.Println("relationships report:")
		for lab, cnt := range relsFoundCount {
			fmt.Println("\tNode type:", lab, " found:", cnt, " merged:", relsMergedCount[lab])
		}
		return nil
	},
}

func init() {
	importCmd.AddCommand(jsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	jsonCmd.Flags().StringVarP(&fPath, "filepath", "f", "", "path to the json file")
	jsonCmd.Flags().BoolVarP(&refreshConstraints, "refresh-constraints", "r", false, "attempt to read constraints direct from the database")
}
