package platform

func CreatePlatform(name string, url string) *PlatformT {
	platform := &PlatformT{name: name, url: url}
	platform.games = make(map[string]string)
	return platform
}

type PlatformT struct {
	name  string
	url   string
	games map[string]string
}

func (p *PlatformT) GetURL() string {
	return p.url
}

func (p *PlatformT) CreateUser(username string, password string) (*UserT, error) {
	user := &UserT{
		username: username,
		password: password,
		platform: p,
	}

	if err := user.initialize(); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PlatformT) AddGame(appid string, name string) {
	p.games[appid] = name
}
