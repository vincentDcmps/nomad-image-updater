package nomadfile

import (
	"fmt"
	"nomad-image-updater/internal/dockerImage"
	"os"
	"path/filepath"
	"slices"
)

type Nomadfile struct {
	Path   string
	Images []*dockerImage.DockerImage
}

func GetNomadFiles(path string) []*Nomadfile {

	nomadextention := []string{".hcl", ".nomad"}
	nomadfiles := []*Nomadfile{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if slices.Contains(nomadextention, filepath.Ext(path)) {

			var nomadfile Nomadfile
			nomadfile.Path = path
			nomadfiles = append(nomadfiles, &nomadfile)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return nomadfiles
}
