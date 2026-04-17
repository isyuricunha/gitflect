package provider

type Repo struct {
	Name     string
	Private  bool
	CloneURL string // URL to clone from (often containing embedded token for pulling)
}

type Source interface {
	ListRepos(visibility string) ([]Repo, error)
}

type Destination interface {
	EnsureRepo(name string, private bool) (pushURL string, err error)
}
