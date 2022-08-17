package pkg

import "testing"

type test struct {
	name string
	cfg  Configuration
	want error
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
		{name: "simple test", cfg: Configuration{Username: "Username"}, want: nil},
		{name: "system reserved", cfg: Configuration{Username: "_Username"}, want: errUsernameMisconfig},
		{name: "restricted character", cfg: Configuration{Username: "Username!"}, want: errUsernameMisconfig},
	}
	for _, tc := range tests {
		got := tc.cfg.ValidateUsername()
		if tc.want != got {
			t.Errorf("%s: wanted %v, but got %v", tc.name, tc.want, got)
		}
	}
}
