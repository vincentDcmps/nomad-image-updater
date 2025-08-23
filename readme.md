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

all setting can be overide by an env variable with a prefix "NID\_"

### settings

#### remoteCustomOption

array containg a map of two following value:

- contain: string to check if option need to be apply on docker repository
- options: possible option are: username,password and insecureTLS

#### LoggerOption

#### Git

#### GetTagReplaceURL

## ToDo

- use cobra for command management
- manage pull request from gitea
- manage url in argument to got directly git forge
- use a meta in task to got release note link
- create test
