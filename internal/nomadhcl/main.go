package nomadhcl

import (
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
)

type NomadFile struct {
	Job Job `hcl:"job,block"`
}

type Job struct {
	Name   string   `hcl:"name,label"`
	Groups []Group  `hcl:"group,block"`
	Other  hcl.Body `hcl:"other,remain"`
}

type Group struct {
	Name  string   `hcl:"name,label"`
	Tasks []Task   `hcl:"task,block"`
	Other hcl.Body `hcl:"other,remain"`
}

type Task struct {
	Name   string   `hcl:"name,label"`
	Driver string   `hcl:"driver"`
	Config Config   `hcl:"config,block"`
	Other  hcl.Body `hcl:"other,remain"`
}

type Config struct {
	Image string   `hcl:"image"`
	Other hcl.Body `hcl:"other,remain"`
}

func ParseNomadFile(path string) *NomadFile {
	parser := hclparse.NewParser()
	f, _ := parser.ParseHCLFile(path)
	var nomadFileInstance NomadFile
	gohcl.DecodeBody(f.Body, nil, &nomadFileInstance)
	return &nomadFileInstance
}

func WriteNomadFile(file *NomadFile, path string) bool {

	hclbuffer := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(file, hclbuffer.Body())
	fileWriter, err := os.Create(path)
	if err != nil {
		return false
	}
	fileWriter.Write(hclbuffer.Bytes())
	return false
}
