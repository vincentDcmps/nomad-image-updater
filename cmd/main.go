package main

import (
	"fmt"
	"nomad-image-updater/internal/dockerImage"
	"nomad-image-updater/internal/nomadfiles"
	"os"
)

func main() {
	target := os.Args[1]
	nomadfiles := nomadfile.GetNomadFiles(target)
	var refImages dockerImage.DockerImageslist
	for _, file := range nomadfiles {
		fileimages := dockerImage.NewDockerImageFromNomadFile(file.Path)
		for _, fileimage := range fileimages {
			imageptr := refImages.Addimage(fileimage)
			file.Images = append(file.Images, imageptr)
		}
	}
	fmt.Println(len(refImages))
	for _, image := range refImages {
		fmt.Println(image)
		image.GetUpdate()
		if image.Update {
			fmt.Println(image)
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
