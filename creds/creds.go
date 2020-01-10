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

// Cr459 is for rave-era bot
var CrRave = Credential{50, 49, 55, 51, 54, 51, 52, 54, 57, 52, 59, 66, 66, 72, 84, 123, 110, 52, 90, 51, 109, 66, 49, 53, 85, 123, 49, 119, 114, 68, 82, 54, 108, 50, 67, 86, 66, 114, 87, 68, 57, 107, 115, 50, 87, 104}

func (cr Credential) String() string {
	unmagic(cr)
	return string(cr)
}
