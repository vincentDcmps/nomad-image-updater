package main

import (
	"fmt"
	"log/slog"
	"nomad-image-updater/internal/config"
	"nomad-image-updater/internal/dockerImage"
	"nomad-image-updater/internal/git"
	"nomad-image-updater/internal/nomadfiles"
	"os"
)

func main() {
	config := config.GetConfig()
	target := os.Args[1]
	lvl := &slog.LevelVar{}
	lvl.UnmarshalText([]byte(config.LoggerOption.Level))
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
	nomadfiles := nomadfile.GetNomadFiles(target)
	var refImages dockerImage.DockerImageslist
	var GitUpdater *git.GitUpdater
	if config.Git.Enabled {
		slog.Info("Git Branch cration is enable")
		var err error
		GitUpdater, err = git.NewGitUpdater(target, config.Git.RefBranch)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
	for _, file := range nomadfiles {
		fileimages := dockerImage.NewDockerImageFromNomadFile(file.Path)
		for _, fileimage := range fileimages {
			imageptr := refImages.Addimage(fileimage)
			file.Images = append(file.Images, imageptr)
		}
	}
	slog.Info(fmt.Sprintf("container image to process: %d", len(refImages)))
	for _, image := range refImages {
		slog.Debug(fmt.Sprintf("proccessing image %s", image.Name))
		image.GetUpdate()
		if image.Update {
			slog.Info("image to update:",
				"name", image.Name,
				"OldVersion", image.Tag,
				"NewVersion", image.NewTag)
		}
	}
	for _, nomadfile := range nomadfiles {
		var gitfileupdater *git.GitFileUpdater
		var err error
		if config.Git.Enabled == true {
			gitfileupdater, err = GitUpdater.NewGitFileUpdater(nomadfile)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

		}
		for _, image := range nomadfile.Images {
			if image.Update {
				image.UpdateNomadFile(nomadfile.Path)
				nomadfile.Updated = true
				if config.Git.Enabled {
					gitfileupdater.CommitImage(image)
				}
			}
		}
	}
}
