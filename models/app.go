package models

type DokkuApp struct {
	Name      string
	GitUrl    string
	GitBranch string
	CreatedAt string
	Status    string
	Details   map[string]string
}

type DokkuAppDetails struct {
	Config map[string]string
	Domain Domain
	Report Report
}

type Domain struct {
	Enabled      bool
	AppVhosts    []string
	GlobalVhosts []string
}

type Report struct {
	Dir    string
	Locked bool
}
