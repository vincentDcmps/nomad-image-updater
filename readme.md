# nomad image updater

Aim of this tool is to update docker image in nomad job.

## Arguments

take as argument target folder to lookup nomad file

## image management

if same image with same tag is use in several file tool only check update once for two file.

## tag management

tool manage version extract version number from tag with following form:

- 0
- 0.0
- 0.0.0

if some character are set in prefix or suffix of version lookup only tag with same suffix and prefix

## Config file

config.yaml can be place to following location:

- ./config.yaml
- ~/.config/nomad-image-updater/
- /etc/nomad-image-updater

config file can have follwing setting

### remoteCustomOption

array containg a map of two following value:
- contain:string to check if option need to be apply on docker repository
- options: possible option are: username,password and insecureTLS

## ToDo

- improve log management
- use cobra for command management
- manage pull request from gitea
- manage url in argument
- create docker container
