package nid

import (
	"log/slog"
	"nomad-image-updater/internal/config"
	"nomad-image-updater/internal/git"
)

func Clean () {
	config := config.GetConfig()
	gitUpdater,err := git.NewGitUpdater(".",config.Git.RefBranch)
	if err !=nil {
		slog.Error(err.Error())
		return
	}
	err = gitUpdater.CleanBranch("nomad-image-updater/")
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Sucess clean nomad-image-updater branch")
	
}
