package pkg

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
)

type test struct {
	name string
	cfg  Configuration
	want error
}

func Test_ReadFromFile(t *testing.T) {
	var cfg Configuration
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	// Search config in home directory with name ".geno" (without extension).
	viper.AddConfigPath(path.Join(home, "go", "geno"))
	viper.SetConfigType("yaml")
	viper.SetConfigName(".geno")

	if err := viper.ReadInConfig(); err != nil {
		t.Error(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		t.Error(err)
	}

	if cfg.Server != "localhost:7687" {
		t.Error("server was not read in correctly. wanted localhost:7687, but got", cfg.Server)
	}
	if cfg.Database != "geno" {
		t.Error("database was not read in correctly. wanted geno, but got", cfg.Database)
	}
	if cfg.User != "geno" {
		t.Error("username was not read in correctly. wanted geno, but got", cfg.User)
	}
}

func Test_ValidateServer(t *testing.T) {
	var tests []test = []test{
		{name: "local server", cfg: Configuration{Server: "localhost:7887"}, want: nil},
		{name: "no port", cfg: Configuration{Server: "localhost"}, want: errServerMisconfig},
	}

	for _, tc := range tests {
		got := tc.cfg.ValidateServer()
		if tc.want != got {
			t.Errorf("%s: wanted %v, but got %v", tc.name, tc.want, got)
		}
	}
}

func Test_ValidateDatabase(t *testing.T) {
	var tests []test = []test{
		{name: "simple test", cfg: Configuration{Database: "database"}, want: nil},
		{name: "system reserved", cfg: Configuration{Database: "_database"}, want: errDatabaseMisconfig},
		{name: "restricted character", cfg: Configuration{Database: "database!"}, want: errDatabaseMisconfig},
		{name: "too long", cfg: Configuration{Database: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"}, want: errDatabaseMisconfig},
		{name: "too short", cfg: Configuration{Database: "ab"}, want: errDatabaseMisconfig},
	}
	for _, tc := range tests {
		got := tc.cfg.ValidateDatabase()
		if tc.want != got {
			t.Errorf("%s: wanted %v, but got %v", tc.name, tc.want, got)
		}
	}
}

func Test_ValidateUsename(t *testing.T) {
	var tests []test = []test{
		{name: "simple test", cfg: Configuration{User: "Username"}, want: nil},
		{name: "system reserved", cfg: Configuration{User: "_Username"}, want: errUsernameMisconfig},
		{name: "restricted character", cfg: Configuration{User: "Username!"}, want: errUsernameMisconfig},
	}
	for _, tc := range tests {
		got := tc.cfg.ValidateUsername()
		if tc.want != got {
			t.Errorf("%s: wanted %v, but got %v", tc.name, tc.want, got)
		}
	}
}
