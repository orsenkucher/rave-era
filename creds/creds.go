package creds

import "io/ioutil"

const filepath = "../key.txt"

// ReadToken reads token from hidden file
func ReadToken() (string, error) {
	token, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

// Credential represens credential for telegram bot
type Credential []byte

// CrRvra is for RAVE'ERA BOT
var CrRvra = Credential{50, 49, 52, 54, 52, 49, 57, 51, 53, 51, 59, 66, 66, 71, 99, 52, 102, 111, 58, 88, 49, 56, 50, 55, 56, 122, 82, 110, 87, 103, 96, 79, 102, 88, 78, 119, 50, 80, 79, 108, 90, 87, 105, 67, 69, 112}

func (cr Credential) String() string {
	unmagic(cr)
	return string(cr)
}
