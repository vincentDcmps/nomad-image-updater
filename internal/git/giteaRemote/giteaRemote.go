package giteaRemote

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"code.gitea.io/sdk/gitea"
)

type giteaRemote struct {
	Client *gitea.Client
	Repo string
	Owner string
}

func NewGiteaRemote (giteaurl string, token string) *giteaRemote {
	p,err := url.Parse(giteaurl)
	if(err!=nil){
		slog.Error("url argument not valid", "stage","gitearemote")
		return nil
	}
	client,err := gitea.NewClient(fmt.Sprintf("%s://%s",p.Scheme,p.Host),gitea.SetToken(token))
	if err != nil{
		slog.Debug(err.Error())
	}
	psplit := strings.Split(p.Path,"/")
	slog.Debug(p.Path)
	slog.Debug(fmt.Sprintf("%#v",psplit))
	if len(psplit) < 2 {
		slog.Error("error when parsing owner and reponame", "stage","gitearemote")
		return nil
	}
	return &giteaRemote{
		Client: client,
		Owner: psplit[1],
		Repo: strings.ReplaceAll(psplit[2],".git", "" ),
	}
}

func (g giteaRemote) CreatePR(head string, base string, title string)error{
	_ ,_ , err :=g.Client.CreatePullRequest(g.Owner,g.Repo,gitea.CreatePullRequestOption{
		Head: head,
		Base: base,
		Title: title,
		Body: "",
	})
	return err
}
