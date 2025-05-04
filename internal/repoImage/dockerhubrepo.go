package repoImage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var DockerhubURL = "hub.docker.com/"


type DockertagJSONResponse struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Prev string `json:"prev"`
	Results []struct{
		Name string `json:"name"`
		Digest string `json:"digest"`
		Id int `json:"id"`
		Last_updated string `json:"last_updated"`
	}
}


type DockerhubRepo struct {
}

func (d *DockerhubRepo) Getreleases(host string,name string) []string{
	if strings.Contains(name,"/") == false{
		name = fmt.Sprintf("library/%s",name)
	}
	if (len(host)==0){
		host = DockerhubURL
	}
	url := fmt.Sprintf("https://%s/v2/repositories/%s/tags?page_size=1000&ordering=last_updated",host,name)
	resp,err := http.Get(url)
	if err != nil {
		 return nil
	}
	defer resp.Body.Close()
	var dockertagResponse DockertagJSONResponse
	json.NewDecoder(resp.Body).Decode(&dockertagResponse)
	var res []string
	for _,result := range dockertagResponse.Results{
		res = append(res,result.Name)
	}
	return res
}

func (d *DockerhubRepo) Validaterepo(repo string ) bool{
	if  (repo == DockerhubURL|| len(repo) == 0)  {
		fmt.Println("docker")
		return true
	}else{
		return false
	}
}


