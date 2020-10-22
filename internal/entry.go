package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	git "github.com/libgit2/git2go/v30"
)

func (app *App) credentialsCallback(url string, username string, _ git.CredType) (*git.Cred, error) {
	return git.NewCredSshKey("git", app.HomeDir+"/.ssh/id_rsa.pub", app.HomeDir+"/.ssh/id_rsa", "")
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

func (app *App) Initialize(repo string) error {
	cloneOptions := &git.CloneOptions{}
	cloneOptions.FetchOptions = &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      app.credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}
	_, err := git.Clone(repo, app.HomeDir+"/.jot", cloneOptions)
	return err
}

func (app *App) Create() error {
	f, err := ioutil.TempFile("/tmp", "entry*.md")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	cmd := exec.Command("emacs", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(f.Name())
	fmt.Println(string(bytes))
	return nil
}
