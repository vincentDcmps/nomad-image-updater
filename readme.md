# nomad image updater

Aim of this tool is to update docker image in nomad job.

## command

_if use with git feature need to be launch from git repository root_

### Update

take as argument target folder or file to lookup nomad file

### Clean

clean all branch in repository beginning by nomad-image-updater/

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

- ./nomad-image-updater.yaml
- ~/.config/nomad-image-updater/
- /etc/nomad-image-updater

all setting can be overide by an env variable with a prefix "NID\_"

### settings

#### remoteCustomOption

array containg  map of two following value:

- contain: string to check if option need to be apply on docker repository
- options: possible option are: username,password and insecureTLS

#### LoggerOption

- verbose

#### Git

- enabled
- refbranch: branch name where new branch will be base

#### GetTagReplaceURL

array with  map toreplace a docker repository URL by another
- target: target url to replace
- replace: replacement URL

## ToDo

- get image update in go routine
- use a meta in task to got release note link
- create test
