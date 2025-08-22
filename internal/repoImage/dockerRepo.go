package repoImage

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"nomad-image-updater/internal/config"
	"strings"

	"github.com/gboddin/go-www-authenticate-parser"
)

type DockerRepo struct {
}

func (d *DockerRepo) Getreleases(host string, name string) []string {
	host = replaceUrl(host)
	tagsurl, _ := url.Parse(fmt.Sprintf("https://%s/v2/%s/tags/list", host, name))
	authHeader, err := getDockerAuth(tagsurl.Host, name)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", tagsurl.String(), nil)
	if authHeader != "" {
		req.Header.Add("authorization", fmt.Sprintf("Bearer %s", authHeader))
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	var tagsListResponse tagsListResponse
	json.NewDecoder(resp.Body).Decode(&tagsListResponse)
	
	slog.Debug(fmt.Sprintf("%d tag found",len(tagsListResponse.Tags)),"image",name)
	return tagsListResponse.Tags

}

func (d *DockerRepo) Validaterepo(repo string) bool {
	if len(repo) > 0 {
		return true
	} else {
		return false
	}
}

func getDockerAuth(host string, name string) (string, error) {
	opt := getRemoteOptions(host)
	if opt.InsecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}
	resp, err := http.Get(fmt.Sprintf("https://%s/v2/_catalog", host))
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		return "", nil
	} else if resp.StatusCode == 401 {
		wwwAuthenticateHeader := www_authenticate_parser.Parse(resp.Header.Get("www-authenticate"))
		realm := wwwAuthenticateHeader.Params["realm"]
		service := wwwAuthenticateHeader.Params["service"]
		scope := fmt.Sprintf("repository:%s:pull", name)
		urlToken := fmt.Sprintf("%s?scope=%s&service=%s", realm, scope, service)
		client := http.Client{}
		req, err := http.NewRequest("GET", urlToken, nil)
		if opt.Username != "" && opt.Password != "" {
			req.Header.Add("Authorisation", basicAuth(opt.Username, opt.Password))
		}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		var tokenResponse tokenResponse
		json.NewDecoder(resp.Body).Decode(&tokenResponse)
		return tokenResponse.Token, nil

	} else {
		return "", errors.New(fmt.Sprintf("unmanage auth return code %d",resp.StatusCode))
	}

}

type tokenResponse struct {
	Token string `json:"token"`
}

type tagsListResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func replaceUrl(host string) string{
	config := config.GetConfig()
	for _,getTagReplaceURL := range config.GetTagReplaceURL {
		if strings.Contains(host,getTagReplaceURL.Target){
			return getTagReplaceURL.Replace
		}
	}
	return host
}

func getRemoteOptions(host string) config.RemoteOptions {
	option := config.RemoteOptions{}
	config := config.GetConfig()
	for _, remote := range config.RemoteCustomOption {
		if strings.Contains(host, remote.Contain) {
			option.Merge(remote.Options)
		}
	}
	return option
}

//"https://ghcr.io/token?scope=repository:docker-mailserver/docker-mailserver:pull&service=ghcr.io" -a "vincent@ducamps.eu:githubpat"
//http https://ghcr.io/v2/docker-mailserver/docker-mailserver/tags/list -A bearer -a Z2hwX1Y2R1ZPcFVGTXdoMkJ3Z3Jtd3FsV0hDc2VFVzJHSjJzMWJsVw==
