package nid

import (
	"fmt"
	"log/slog"
	"nomad-image-updater/internal/config"
	"nomad-image-updater/internal/dockerImage"
	"nomad-image-updater/internal/git"
	"nomad-image-updater/internal/git/giteaRemote"
	"nomad-image-updater/internal/nomadfiles"
)

func Update(target string) {
	config := config.GetConfig()
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
		if (config.Git.RemoteCreatePR == "gitea"){
			GitUpdater.Remote= giteaRemote.NewGiteaRemote(config.Git.RemoteURL,config.Git.RemoteToken)
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
			slog.Info("image to update",
				"name", image.Name,
				"OldVersion", image.Tag,
				"NewVersion", image.NewTag)
		}
	}
	for _, nomadfile := range nomadfiles {
		var gitfileupdater *git.GitFileUpdater
		var err error
		for _, image := range nomadfile.Images {
			if image.Update {
				if gitfileupdater == nil && config.Git.Enabled == true {
					gitfileupdater, err = GitUpdater.NewGitFileUpdater(nomadfile)
					if err != nil {
						slog.Error(err.Error())
						break
					}
				}
				if config.Git.Enabled {
					gitfileupdater.CommitImage(image)
				}else{
					image.UpdateNomadFile(nomadfile.Path)
				}
				nomadfile.Updated = true
			}
		}

		if (nomadfile.Updated ==true && config.Git.Enabled ==true && config.Git.RemoteURL != "" && config.Git.RemoteToken != ""){
			err:=gitfileupdater.Push(config.Git.RemoteURL,config.Git.RemoteToken)
			if(err != nil){
				slog.Error(err.Error())
				break
			}
			slog.Debug(fmt.Sprintf("%#v",gitfileupdater.GitUpdater.Remote))
			if gitfileupdater.GitUpdater.Remote != nil {
				gitfileupdater.CreatePR()
			}
		}
	}
}
