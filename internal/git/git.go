package git

import (
	"fmt"
	"log/slog"
	"nomad-image-updater/internal/dockerImage"
	"nomad-image-updater/internal/nomadfiles"
	"os"

	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"strings"
)

type GitRemoteInterface interface {
	CreatePR(head string, base string, title string) error
}

type GitUpdater struct {
	Repository      *git.Repository
	ReferenceBranch plumbing.ReferenceName
	Remote          GitRemoteInterface
}

func NewGitUpdater(target string, refname string) (*GitUpdater, error) {
	file, _ := os.Open(target)
	defer file.Close()
	fileinfo, _ := file.Stat()
	if fileinfo.IsDir() {
		r, err := git.PlainOpen(target)
		if err == nil {
			slog.Debug(fmt.Sprintf("Repo found in target: %s", target))
			return &GitUpdater{Repository: r, ReferenceBranch: plumbing.ReferenceName(refname)}, nil
		}
	}
	currentpath, _ := os.Getwd()
	r, err := git.PlainOpen(currentpath)
	if err != nil {
		slog.Debug("repository not found")
		return &GitUpdater{}, err
	}
	slog.Debug(fmt.Sprintf("Repo found in %s", currentpath))
	refbranchConfig, err := r.Branch(refname)
	if err != nil {
		slog.Error(err.Error())
		return &GitUpdater{}, err
	}
	return &GitUpdater{Repository: r, ReferenceBranch: refbranchConfig.Merge}, nil
}

func (g *GitUpdater) NewGitFileUpdater(file *nomadfile.Nomadfile) (*GitFileUpdater, error) {
	slog.Debug(file.GetFileName())
	baseref, _ := g.Repository.Storer.Reference(g.ReferenceBranch)
	localref := plumbing.ReferenceName(fmt.Sprintf("refs/heads/nomad-image-updater/%s", file.GetFileName()))
	var c *plumbing.Reference
	g.Repository.Head()
	c = plumbing.NewHashReference(localref, baseref.Hash())
	g.Repository.Storer.SetReference(c)
	return &GitFileUpdater{
		File:       file,
		GitUpdater: g,
		Branch:     c,
	}, nil
}

func (g *GitUpdater) CleanBranch(pattern string) error {
	branches, err := g.Repository.Branches()
	if err != nil {
		return err
	}
	err = branches.ForEach(func(r *plumbing.Reference) error {
		name := r.Name()
		if strings.Contains(name.String(), pattern) {
			slog.Info(fmt.Sprintf("deleting branch %s", name.String()))
			err := g.Repository.Storer.RemoveReference(name)
			return err
		}
		return nil
	})
	return err
}

type GitFileUpdater struct {
	File       *nomadfile.Nomadfile
	GitUpdater *GitUpdater
	Branch     *plumbing.Reference
}

func (g *GitFileUpdater) CommitImage(image *dockerImage.DockerImage) bool {
	slog.Debug(fmt.Sprintf("commit of image %s in file %s", g.File.Path, image.Name))
	w, err := g.GitUpdater.Repository.Worktree()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: g.Branch.Name(),
	})
	if err != nil {
		slog.Error(err.Error(), "stage", "Checkout")
		return false
	}

	image.UpdateNomadFile(g.File.Path)
	_, err = w.Add(g.File.Path)
	if err != nil {
		slog.Error(err.Error(), "stage", "add")
		return false
	}
	commitmsg := fmt.Sprintf("update (%s): %s to %s", g.File.GetFileName(), image.Name, image.NewTag)
	_, err = w.Commit(commitmsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "nomad-image-updater",
			Email: "",
			When:  time.Now(),
		},
	})
	if err != nil {
		slog.Error(err.Error(), "stage", "commit")
		return false
	}
	w.Checkout(&git.CheckoutOptions{
		Branch: g.GitUpdater.ReferenceBranch,
	})

	return true
}

func (g *GitFileUpdater) Push(remote string, token string) error {
	refspec := config.RefSpec(fmt.Sprintf("%s:%s", g.Branch.Name().String(), g.Branch.Name().String()))
	slog.Debug(refspec.String())
	pushoptions := git.PushOptions{
		RemoteURL: remote,
		RefSpecs:  []config.RefSpec{refspec},
		Force:     true,
		Auth: &http.TokenAuth{
			Token: token,
		}}
	err := g.GitUpdater.Repository.Push(&pushoptions)
	return err
}

func (g *GitFileUpdater) CreatePR() {
	slog.Info(fmt.Sprintf("create PR for %s", g.Branch.Name().Short()))
	err := g.GitUpdater.Remote.CreatePR(g.Branch.Name().Short(), g.GitUpdater.ReferenceBranch.Short(), g.Branch.Name().Short())
	if err != nil {
		slog.Error(err.Error())
	}
}
