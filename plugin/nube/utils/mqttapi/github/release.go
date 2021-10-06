package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Release struct {
	ID         uint   `json:"id"`
	ZipBallURL string `json:"zipball_url"`
	TagName    string `json:"tag_name"`
}

// API interface
type API interface {
	GetLatestReleaseTag(string) (string, error)
	GetReleaseTags(string) ([]string, error)
}

// GitHub GitHub struct
type GitHub struct{}

func New() *GitHub {
	return &GitHub{}
}

// GetLatestReleaseTag returns the latest release tag
func (gh *GitHub) GetLatestReleaseTag(repo string) (string, error) {
	tags, err := gh.GetReleaseTags(repo)
	if err != nil {
		fmt.Println("error: GetLatestReleaseTag")
		return "", err
	}
	return tags[0], nil
}

// GetReleaseTags returns all releases tag name per a repository
func (gh *GitHub) GetReleaseTags(repo string) ([]string, error) {
	var tags []string
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	fmt.Println("apiURL:", apiURL)
	res, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("error: GetReleaseTags")
		return tags, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error: GetReleaseTags")
		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error: GetReleaseTags")
		return tags, err
	}
	var releases []Release
	if err = json.Unmarshal(body, &releases); err != nil {
		fmt.Println("error: GetReleaseTags")
		return tags, err
	}
	tags = make([]string, len(releases))
	for i, release := range releases {
		tags[i] = release.TagName
	}

	return tags, nil
}

func GetRelease(repo string, version string, token string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	if version != "latest" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", repo, version)
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	if err != nil {
		return nil, err
	}
	client := http.Client{Timeout: 30 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if err = checkResponse(r); err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(r.Body)
	var release Release
	if err = json.NewDecoder(r.Body).Decode(&release); err != nil {
		return nil, err
	}
	if release.TagName == "" || release.ZipBallURL == "" {
		return nil, fmt.Errorf("release info missing tag_name and/or zipball_url")
	}

	if version == "latest" {
		fmt.Println("version", release.TagName)
	}
	return &release, nil
}

func (release *Release) Download() (string, error) {

	tmpFile, err := ioutil.TempFile("", "nube-download-")
	if err != nil {
		return "", err
	}

	rz, err := http.Get(release.ZipBallURL)
	if err != nil {
		return "", err
	}
	if err = checkResponse(rz); err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(rz.Body)

	if _, err = io.Copy(tmpFile, rz.Body); err != nil {
		return "", err
	}

	if err = tmpFile.Sync(); err != nil {
		return "", err
	}
	if err = tmpFile.Close(); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("HTTP 404: %s", r.Request.URL)
	}

	var e struct{ Message string }

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		if err = json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("could not parse error HTTP %d error: %s", r.StatusCode, data)
		}
	}

	if e.Message != "" {
		return fmt.Errorf("github HTTP %d error: %s", r.StatusCode, e.Message)
	}

	return fmt.Errorf("github HTTP %d error: %s", r.StatusCode, data)
}

//func main() {
//
//a := New()
//tags, err := a.GetReleaseTags("nube/server")
//if err != nil {
//	return

//
//	tag, err := a.GetReleaseTags("nube/server")
//	if err != nil {
//		return
//	}
//	fmt.Println(tag)
//	tags, err := a.GetLatestReleaseTag("nube/server")
//	if err != nil {
//		return
//	}
//	fmt.Println(tags)
//
//
//	release, err := GetRelease("aaa", "v2.0.23", "ghp_j9mu")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(release)
//	download, err := release.Download()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	uz := unzip.New()
//
//	files, err := uz.Extract("/tmp/youtube-dl-850842546", "./data/directory")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(files)
//	fmt.Println(download)
//
//}
