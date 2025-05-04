package dockerImage

import (
	"fmt"
	"nomad-image-updater/internal/repoImage"
	"os"
	"regexp"
	"github.com/hashicorp/go-version"
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
	"version1": `^(?P<prefix>[1-9A-Za-z\-]*)(?P<version>(\d+))(?P<suffix>[1-9A-Za-z\-]*)$`,
	"version2": `^(?P<prefix>[1-9A-Za-z\-]*)(?P<version>(\d+\.)(\*|\d+))(?P<suffix>[1-9A-Za-z\-]*)$`,
	"version3": `^(?P<prefix>[1-9A-Za-z\-]*)(?P<version>(\d+\.)(\d+\.)(\*|\d+))(?P<suffix>[1-9A-Za-z\-]*)$`,
	"latest": `latest|^$`,
}

type DockerImage struct {
	URL    string
	Name   string
	Tag    string 
	TagType string
	NewTag string
	Update bool
}

func NewDockerImage(URL string, name string, tag string) DockerImage {
	var im DockerImage
	im.URL = URL
	im.Name = name
	im.Tag = tag
	im.Update = false
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
	imageRegex, _ := regexp.Compile(`image\s*=\s*\"(?P<repo>(?P<URL>[^:@\n\/]*\.[^:@\n\/]*(:\d*)?\/)?(?P<image>[^:@\n]*))(:(?P<tag>[^:@\n]*))?(@.*:.*)?\"`)
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
		resp = append(resp, &image)
	}
	return resp
}
 
func(d * DockerImage) ToString(newtag bool) string{
	if newtag && d.NewTag != ""{
		return fmt.Sprintf("%s%s:%s",d.URL,d.Name,d.NewTag)
	}else{
		return fmt.Sprintf("%s%s:%s",d.URL,d.Name,d.Tag)
	}
}


func(d *DockerImage) UpdateNomadFile(path string){
	f,err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	r,_ := regexp.Compile(d.ToString(false))
	fmt.Println(d.ToString(false))
	fmt.Println(d.ToString(true))
	newfile := r.ReplaceAllString(string(f),d.ToString(true))
  stat, err :=	os.Stat(path)
	err = os.WriteFile(path,[]byte(newfile),stat.Mode())
	
	
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
	if d.TagType == "latest"{
		return 
	}else if(len(d.TagType) == 0){
		fmt.Println("No tag type detected")
		return
	}
	r ,_ := regexp.Compile(tagtype[d.TagType])
	var suffix string
	var prefix string
	var lastversion *version.Version
	match := r.FindStringSubmatch(d.Tag)
	suffix = match[r.SubexpIndex("suffix")]
	prefix = match[r.SubexpIndex("prefix")]
	lastversion, _ = version.NewVersion(match[r.SubexpIndex("version")])
	taglist := d.getTags()
	for _,tag := range(taglist) {
		match2 := r.FindStringSubmatch(tag)
		if len(match2) > 0 && match2[r.SubexpIndex("suffix")] == suffix && match2[r.SubexpIndex("prefix")] == prefix {
			version, _ := version.NewVersion(match2[r.SubexpIndex("version")])
			if version.GreaterThan(lastversion){
				lastversion = version
				d.NewTag=tag
				d.Update=true
			}
		}
		
	}
}
