package main

import (
	"fmt"
	"nomad-image-updater/internal/nomadhcl"
	"os"
	"path/filepath"
	"slices"
	"github.com/hashicorp/nomad/jobspec2"
)

func main() {
	target := os.Args[1]
	_, err := os.Stat(target);

	if  os.IsNotExist(err){
		fmt.Println("target file or folder not Existing: ", target)
		os.Exit(-1)	
	}
	nomadfile := getNomadFile(target)

	for _,file := range(nomadfile) {

		parsedhcl := nomadhcl.ParseNomadFile(file)
		if parsedhcl.Job.Name != ""  {
			fmt.Println(parsedhcl.Job.Groups[0].Tasks[0].Config.Image)
		}
		//test jobspec2
		f,_ := os.Open(file)
		job,_ := nomad.jobspec2.Parse(file,f)
		fmt.Println(job.TaskGroups[0].Tasks[0].Config)

	}
}

func getNomadFile (path string) []string {
	
	nomadextention := []string{".hcl",".nomad"}
	nomadfiles := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err!=nil {
					return err
				}
				if slices.Contains(nomadextention,filepath.Ext(path)){
      		nomadfiles=append(nomadfiles,path)
				}
				return nil
			})
		if err != nil {
			fmt.Println(err)
		}

	return nomadfiles
}
