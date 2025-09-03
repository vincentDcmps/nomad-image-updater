package repoImage

import (
	"nomad-image-updater/internal/config"
)

type repoImage interface {
	Getreleases(host string, name string, remoteOption config.RemoteOptions) []string
	Validaterepo(string) bool
}

func GetMapRepo() map[string]repoImage {
	m := map[string]repoImage{
		"dockerhub": &DockerhubRepo{},
		"docker":    &DockerRepo{},
	}
	return m
}
