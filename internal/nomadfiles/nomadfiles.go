package nomadfile

import (
	"fmt"
	"nomad-image-updater/internal/dockerImage"
	"os"
	"path/filepath"
	"slices"
)

type Nomadfile struct {
	Path    string
	Images  []*dockerImage.DockerImage
	Updated bool
}

func GetNomadFiles(path string) []*Nomadfile {

	nomadextention := []string{".hcl", ".nomad"}
	nomadfiles := []*Nomadfile{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if slices.Contains(nomadextention, filepath.Ext(path)) {

			nomadfile := Nomadfile{
				Path:    path,
				Updated: false,
			}
			nomadfiles = append(nomadfiles, &nomadfile)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return nomadfiles
}

func (n *Nomadfile) GetFileName() string {
	return filepath.Base(n.Path)

}
