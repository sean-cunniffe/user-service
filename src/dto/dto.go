package dto

import "time"

// User holds data about user that can be used to populate server responses, and populate
// data storage
type User struct {
	ID        string
	Name      string
	Email     string
	Password  Password
	Role      Role
	Created   time.Time
	LastLogin time.Time
}

// Password contains the hash of the user password and the algorithm used to encode it
type Password struct {
	Hash []byte
	Alg  string
}

// Role contains the permissions of a user
type Role struct {
	Name string
	Perm []string
}

type TLSConfig struct {
	SkipVerify bool   `yaml:"skipVerify" json:"skipVerify"`
	CertFile   string `yaml:"certFile" json:"certFile"`
	KeyFile    string `yaml:"keyFile" json:"keyFile"`
	CAFile     string `yaml:"caFile" json:"caFile"`
}
