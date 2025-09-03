package repoImage

import (
	"nomad-image-updater/internal/config"
	"net/http"
	"crypto/tls"
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


func httpclient (insecure bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	return &http.Client{Transport: tr}
}
