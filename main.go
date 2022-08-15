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
package main

import "github.com/Viking2012/geno/cmd"

// type TestConfiguration struct {
// 	Server      string                      `json:"server"`
// 	Database    string                      `json:"database"`
// 	Username    string                      `json:"username"`
// 	Constraints map[string]geno.Constraints `json:"constraints"`
// }

func main() {
	// as a CLI, this will ultimately be the only command within the application that runs
	cmd.Execute()

	// raw, err := ioutil.ReadFile(path.Join(".", "test", "test_config.json"))
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = pkg.ConfigFromBytes(raw)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("server:", cfg.Server)
	// fmt.Println("database:", cfg.Database)
	// for db, cons := range cfg.Constraints {
	// 	fmt.Println(db)
	// 	fmt.Println("NODE-UNIQUENESS")
	// 	for _, c := range cons.NodeUniqueness {
	// 		fmt.Println(c)
	// 	}
	// 	fmt.Println("NODE-KEYS")
	// 	for _, c := range cons.NodeKeys {
	// 		fmt.Println(c)
	// 	}
	// 	fmt.Println("NODE-PROP-EXISTENCE")
	// 	for _, c := range cons.NodePropertyExistence {
	// 		fmt.Println(c)
	// 	}
	// 	fmt.Println("Relationship-UNIQUENESS")
	// 	for _, c := range cons.RelationshipUniqueness {
	// 		fmt.Println(c)
	// 	}
	// 	fmt.Println("Relationship-KEYS")
	// 	for _, c := range cons.RelationshipKeys {
	// 		fmt.Println(c)
	// 	}
	// 	fmt.Println("Relationship-PROP-EXISTENCE")
	// 	for _, c := range cons.RelationshipPropertyExistence {
	// 		fmt.Println(c)
	// 	}
	// }
}
