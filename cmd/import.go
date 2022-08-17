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
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Insert data into a neo4j database",
	Long: `A collaction of commands which insert data (nodes and relationships)
into a neo4j database from a variety of filetypes. Insertion of data can respect
constraints which extended beyond those in native neo4j.

Filetypes currently included are:
- json (command json)

Non-native constraint types include:
- Uniqueness of nodes with multiple property definitions in Community Edition
- Node keys in Community Edition
- Node property existence on multiple properties for each node label in Community Edition
- Relationship unqiueness
- Relationship keys
- Relationship property existence in Community Edition`,
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	importCmd.PersistentFlags().StringVarP(&cfg.Database, "database", "d", cfg.Database, "Insert records into this database")

	importCmd.PersistentFlags().StringVarP(&cfg.Server, "server", "s", cfg.Server, "Location of database in format: <SERVER>:<PORT>")

	importCmd.PersistentFlags().StringVarP(&cfg.User, "username", "u", cfg.User, "Username used to connect to the server")
}
