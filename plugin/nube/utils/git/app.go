package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/utils/git/github"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/NubeDev/flow-framework/utils/unzip"
	"strings"
)

func download(token string) error {
	owner := "NubeIO"
	repo := "flow-framework"
	tag := "latest"
	deviceType := "armv7"
	dir := "/"
	pluginsDir := "/data/flow-framework/data/plugins"

	a := github.New()
	var version string
	if tag == "latest" {
		_repo := fmt.Sprintf("%s/%s", owner, repo)
		fmt.Println("repo", _repo)
		tags, err := a.GetLatestReleaseTag(_repo)
		if err != nil {
			fmt.Println("error:", _repo)
		}
		version = tags
	} else {
		version = tag
	}
	fmt.Println("version:", version)
	downloads, err := github.RetrieveAssets(owner, repo, version, token, deviceType)
	if err != nil {
		fmt.Println("error: main github.RetrieveAssets:", err)
		return err
	}
	fmt.Println(dir)
	fmt.Println(downloads)
	uz := unzip.New()

	err = utils.DirRemoveContents(pluginsDir)
	if err != nil {
		return errors.New("error on delete all plugins")
	}
	for _, v := range downloads.Values() {
		n := v.(string)
		if !strings.Contains(n, "flow-framework") { //dont unzip main app build
			_, err = uz.Extract(n, pluginsDir)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}
