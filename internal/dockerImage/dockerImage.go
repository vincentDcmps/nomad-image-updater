package dockerImage

import (
	"fmt"
	"nomad-image-updater/internal/repoImage"
	"os"
	"regexp"
)

type DockerImageslist []*DockerImage

func (i *DockerImageslist) Foundimage(lookupimage *DockerImage) *DockerImage {
	for _, image := range *i {
		if *image == *lookupimage {
			return image
		}
	}
	return nil
}
func (i *DockerImageslist) Addimage(image *DockerImage) *DockerImage {
	ptrimage := i.Foundimage(image)
	if ptrimage == nil {
		*i = append(*i, image)
		ptrimage = image
	}
	return ptrimage

}

var tagtype = map[string]string{
	"version": `(?P<prefix>[1-9A-Za-z\-]*)(?P<version>(\d+\.)?(\d+\.)?(\*|\d+))(?P<suffix>[1-9A-Za-z\-]*)`,
	"latest": "latest|",
}

type DockerImage struct {
	URL    string
	Name   string
	Tag    string 
	TagType string
	update bool
}

func NewDockerImage(URL string, name string, tag string) DockerImage {
	var im DockerImage
	im.URL = URL
	im.Name = name
	im.Tag = tag
	im.update = false
	for k,v :=range tagtype {
		tagtypeRegex,_ := regexp.Compile(v)
		match := tagtypeRegex.MatchString(im.Tag)
		if match {
			im.TagType = k
			break
		}
	}
	return im
}

func NewDockerImageFromNomadFile(path string) DockerImageslist {
	var resp DockerImageslist
	imageRegex, _ := regexp.Compile(`image\s*=\s*\"(?P<repo>(?P<URL>[^:@\n]*(:\d*)\/)?(?P<image>[^:@\n]*))(:(?P<tag>[^:@\n]*))?(@.*:.*)?\"`)
	f, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	matches := imageRegex.FindAllStringSubmatch(string(f), -1)
	for _, match := range matches {
		image := NewDockerImage(match[imageRegex.SubexpIndex("URL")],
			match[imageRegex.SubexpIndex("image")],
			match[imageRegex.SubexpIndex("tag")])
		if image.URL == "" {
			image.URL = repoImage.DockerhubURL
		}
		resp = append(resp, &image)
	}
	return resp
}

func (d *DockerImage) getTags() []string {
	for _, v := range repoImage.GetMapRepo() {
		if v.Validaterepo(d.URL) {
			releasesList := v.Getreleases(d.URL, d.Name)
			return releasesList
		}
	}
	return nil
}


func (d *DockerImage) GetUpdate() {
	taglist := d.getTags()
	r ,_ := regexp.Compile(tagtype[d.TagType])
	var suffix string
	var prefix string
	var filteringTag []string
	if d.TagType == "latest" {
		return 
	}else{
		match := r.FindStringSubmatch(d.Tag)
		suffix = match[r.SubexpIndex("suffix")]
		prefix = match[r.SubexpIndex("prefix")]
	}
	for _,tag := range(taglist) {
		match2 := r.FindStringSubmatch(tag)
		if len(match2) > 0 && match2[r.SubexpIndex("suffix")] == suffix && match2[r.SubexpIndex("prefix")] == prefix {
			filteringTag = append(filteringTag, tag )
		}
	}
	fmt.Println(filteringTag)
}
