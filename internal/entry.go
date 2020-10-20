package internal

import (
	git "github.com/libgit2/git2go/v30"
)

func (app *App) credentialsCallback(url string, username string, _ git.CredType) (*git.Cred, error) {
	return git.NewCredSshKey("git", app.HomeDir+".ssh/id_rsa.pub", app.HomeDir+".ssh/id_rsa", "")
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
