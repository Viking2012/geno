package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"syscall"

	"github.com/Viking2012/geno/geno"
	"golang.org/x/term"
)

var (
	errServerMisconfig   error = errors.New("configured server must ccontain an ip address and a port number, seperated by a color")
	errDatabaseMisconfig error = errors.New("configured database must obey neo4j naming conventions")
	errUsernameMisconfig error = errors.New("configured username must obey neo4j naming conventions")
)

type Configuration struct {
	Server      string `mapstructure:"server"`
	Database    string `mapstructure:"database"`
	User        string `mapstructure:"user"`
	password    string
	Constraints map[string]geno.Constraints `mapstructure:"constraints"`
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
	if cfg.Server == "" {
		return errServerMisconfig
	}
	matched, _ := regexp.MatchString(`\:`, cfg.Server)
	if !matched {
		return errServerMisconfig
	}
	return nil
}

// ValidateDatabase checks for all of the rules outlined in:
// https://neo4j.com/docs/operations-manual/current/manage-databases/configuration/
func (cfg *Configuration) ValidateDatabase() error {
	var (
		matched bool
		err     error
	)
	if len(cfg.Database) > 63 {
		return errDatabaseMisconfig
	}
	if len(cfg.Database) < 3 {
		return errDatabaseMisconfig
	}
	matched, err = regexp.MatchString(`^[[:alnum:]]`, cfg.Database) // should match
	if err != nil || !matched {
		return errDatabaseMisconfig
	}
	matched, err = regexp.MatchString(`[^a-zA-Z0-9.-]`, cfg.Database) // should not match
	if err != nil || matched {
		return errDatabaseMisconfig
	}
	matched, err = regexp.MatchString(`^_`, cfg.Database) // should not match
	if err != nil || matched {
		return errDatabaseMisconfig
	}
	return nil
}
func (cfg *Configuration) ValidateUsername() error {
	var (
		matched bool
		err     error
	)
	matched, err = regexp.MatchString(`^[[:alnum:]]`, cfg.User) // should match
	if err != nil || !matched {
		return errUsernameMisconfig
	}
	matched, err = regexp.MatchString(`[^a-zA-Z0-9.-]`, cfg.User) // should not match
	if err != nil || matched {
		return errUsernameMisconfig
	}
	matched, err = regexp.MatchString(`^_`, cfg.User) // should not match
	if err != nil || matched {
		return errUsernameMisconfig
	}
	return nil
}
func (cfg *Configuration) SetPassword(p string) { cfg.password = p }
func (cfg *Configuration) GetPassword() string  { return cfg.password }

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

func (cfg *Configuration) ValidateWithAttempts() error {
	var triedToGetUsername bool = false
	var triedToGetPassword bool = false
	if cfg.Database == "" {
		return errors.New("database name must be provided either via a configuration file (--config) or via the database flag (-d, --database)")
	}
	if cfg.Server == "" {
		return errors.New("server location must be provided either via a configuration file (--config) or via the server flag (-s, --server)")
	}
	for cfg.User == "" {
		if !triedToGetUsername {
			triedToGetUsername = true
			fmt.Println("Username cannot be blank. Set it below, via a configuration file (--config) or via the username flag (-u, --uesrname)")
			err := askUsername(cfg)
			if err != nil {
				fmt.Println("Username could not be set!")
				askUsername(cfg)
			}
		} else {
			return errors.New("username could not be set and cannot be blank. Set it via a configuration file (--config) or via the username flag (-u, --uesrname)")
		}
	}
	for cfg.GetPassword() == "" {
		if !triedToGetPassword {
			triedToGetPassword = true
			fmt.Println("Password cannot be blank. Set it below")
			err := askPassword(cfg)
			if err != nil {
				fmt.Println("Password could not be set!")
				askPassword(cfg)
			}
		} else {
			return errors.New("password could not be set and cannot be blank")
		}
	}
	return nil
}

func askUsername(cfg *Configuration) error {
	fmt.Print("Enter Username: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	cfg.User = string(bytePassword)
	fmt.Println()
	return nil
}

func askPassword(cfg *Configuration) error {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	cfg.SetPassword(string(bytePassword))
	fmt.Println()
	return nil
}
