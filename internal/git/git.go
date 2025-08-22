package git

import (
	"fmt"
	"log/slog"
	"nomad-image-updater/internal/nomadfiles"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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

func (g *GitUpdater) CreateUpdateBranch(file *nomadfile.Nomadfile) bool {
	slog.Debug(file.GetFileName())
	localref := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", file.GetFileName()))
	w, err := g.Repository.Worktree()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: localref, Create: true, Keep: true,
	})
	if err != nil {
		slog.Error(err.Error(), "stage", "Checkout")
		return false
	}
	_, err = w.Add(file.Path)
	if err != nil {
		slog.Error(err.Error(), "stage", "add")
		return false
	}
	_, err = w.Commit(fmt.Sprintf("update: %s", file.GetFileName()), &git.CommitOptions{})
	if err != nil {
		slog.Error(err.Error(), "stage", "commit")
		return false
	}
	w.Checkout(&git.CheckoutOptions{
		Branch: g.ReferenceBranch,
	})
	return true
}
