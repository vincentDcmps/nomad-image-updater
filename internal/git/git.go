package git

import (
	"fmt"
	"log/slog"
	"nomad-image-updater/internal/dockerImage"
	"nomad-image-updater/internal/nomadfiles"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"time"
)

type GitUpdater struct {
	Repository      *git.Repository
	ReferenceBranch plumbing.ReferenceName
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
		Branch: g.Branch.Name(), Keep: true,
	})
	if err != nil {
		slog.Error(err.Error(), "stage", "Checkout")
		return false
	}
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
