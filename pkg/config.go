package pkg

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/Viking2012/geno/geno"
)

type Configuration struct {
	Server      string
	Database    string
	Username    string
	password    string
	Constraints map[string]geno.Constraints
}

func ConfigFromBytes(raw []byte) (cfg Configuration, err error) {
	err = json.Unmarshal(raw, &cfg)
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
	return cfg, err
}

// ValidateServer ensures that a server and port have been provided
func (cfg *Configuration) ValidateServer() error {
	var err error = errors.New("configured server must ccontain an ip address and a port number, seperated by a color")

	if cfg.Server == "" {
		return err
	}
	matched, _ := regexp.MatchString(":", cfg.Server)
	if !matched {
		return err
	}
	return nil
}
func (cfg *Configuration) ValidateDatabase() error { return nil }
func (cfg *Configuration) ValidateUsername() error { return nil }
func (cfg *Configuration) SetPassword(p string)    { cfg.password = p }
func (cfg *Configuration) GetPassword() string     { return cfg.password }

// Validate ensures that all required configuration values are present and in the correct format
func (cfg *Configuration) Validate() error {
	var err error
	err = cfg.ValidateServer()
	if err != nil {
		return err
	}
	err = cfg.ValidateDatabase()
	if err != nil {
		return err
	}
	err = cfg.ValidateUsername()
	if err != nil {
		return err
	}
	return nil
}
