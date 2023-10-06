package scripts

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadScriptFromGitHub(repo, branch, path, script string) string {
	githubToken := os.Getenv("ODX_GITHUB_TOKEN")
	cacheDir := getCacheDir()
	// TODO: handle same name from multiple sources e.g. forks
	localScript, err := os.Create(cacheDir + string(os.PathSeparator) + script)
	if err != nil {
		return ""
	}

	err = localScript.Chmod(os.FileMode(0755))
	if err != nil {
		return ""
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", repo, branch, path, script), nil)
	request.Header.Set("Accept", "application/vnd.github.v3.raw")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
	response, err := client.Do(request)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err != nil {
		return ""
	}
	_, err = io.Copy(localScript, response.Body)
	if err != nil {
		return ""
	}

	return localScript.Name()
}
