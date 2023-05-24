package repl

// Account is cli's account
type Account struct {
	Username       string                       `json:"Username"`
	Aliases        map[string][]string          `json:"Aliases"`
	CmdSets        map[string]map[string]string `json:"CmdSets"`
	CurrentSet     string                       `json:"CurrentSet"`
	CurrentSetType int                          `json:"CurrentSetType"`
}

// Load loads account
func (a *Account) Load() (bool, error) {

	logger.Printf("%s logined successfully\n", a.Username)
	return false, nil
}

// Save saves account
func (a *Account) Save() error {

	return nil
}
