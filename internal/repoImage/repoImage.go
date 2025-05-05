package repoImage

type repoImage interface {
	Getreleases(host string, name string) []string
	Validaterepo(string) bool
}

func GetMapRepo() map[string]repoImage {
	m := map[string]repoImage{
		"dockerhub": &DockerhubRepo{},
		"docker":    &DockerRepo{},
	}
	return m
}
