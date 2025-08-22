package main

import (
	"fmt"
	"log/slog"
	"nomad-image-updater/internal/dockerImage"
	"nomad-image-updater/internal/nomadfiles"
	"nomad-image-updater/internal/config"
	"os"
)

func main() {
	config := config.GetConfig()
	target := os.Args[1]
	lvl:=&slog.LevelVar{}
	lvl.UnmarshalText([]byte(config.LoggerOption.Level))
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
        Level: lvl,
    }))
	slog.SetDefault(logger)
	nomadfiles := nomadfile.GetNomadFiles(target)
	var refImages dockerImage.DockerImageslist
	for _, file := range nomadfiles {
		fileimages := dockerImage.NewDockerImageFromNomadFile(file.Path)
		for _, fileimage := range fileimages {
			imageptr := refImages.Addimage(fileimage)
			file.Images = append(file.Images, imageptr)
		}
	}
	slog.Info(fmt.Sprintf("container image to process: %d",len(refImages)))
	for _, image := range refImages {
		slog.Debug(fmt.Sprintf("proccessing image %s",image.Name))
		image.GetUpdate()
		if image.Update {
		slog.Info("image to update:",
							"name",image.Name,
							"OldVersion",image.Tag,
							"NewVersion",image.NewTag)
		}
	}
	for _, nomadfile := range nomadfiles {
		for _, image := range nomadfile.Images {
			if image.Update {
				image.UpdateNomadFile(nomadfile.Path)
			}
		}
	}
}
