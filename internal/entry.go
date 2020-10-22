package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v30"
)

var (
	endOfContent = []byte("\n")
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
	_, err := git.Clone(repo, app.GitDir, cloneOptions)
	return err
}

func (app *App) Create() error {
	buf := make([]byte, 0)
	conbuf := bytes.NewBuffer(buf)

	now := time.Now()
	fname := fmt.Sprintf("%s/%s.md", app.GitDir, now.Format("2006-01-02"))

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		dateHeader := now.Format("Monday, January _2, 2006")
		dateHeader = fmt.Sprintf("# %s\n\n", dateHeader)
		conbuf.WriteString(dateHeader)
	}

	timeHeader := now.Format("15:04:05 MST")
	timeHeader = fmt.Sprintf("## %s\n\n", timeHeader)
	conbuf.WriteString(timeHeader)

	content, err := app.getInput()
	if err != nil {
		return err
	}
	conbuf.Write(content)

	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	conbuf.Write(endOfContent)
	if _, err := f.Write(conbuf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (app *App) Update(date string) error {

	if !strings.ContainsAny(date, ".md") { //support not adding extension
		date = date + ".md"
	}

	tf, err := ioutil.TempFile("/tmp", "entry*.md")
	if err != nil {
		return err
	}
	defer tf.Close()
	defer os.Remove(tf.Name())

	fname := fmt.Sprintf("%s/%s", app.GitDir, date)
	exFile, err := os.Open(fname)
	if err != nil {
		return err
	}

	if _, err := io.Copy(tf, exFile); err != nil {
		return err
	}

	if err := tf.Sync(); err != nil {
		return err
	}

	content, err := app.getInputWithExistingFile(tf)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) //OS truncate
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return err
	}

	return nil
}

func (app *App) getInput() ([]byte, error) {
	f, err := ioutil.TempFile("/tmp", "entry*.md")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	return app.getInputWithExistingFile(f)
}

func (app *App) getInputWithExistingFile(f *os.File) ([]byte, error) {
	cmd := exec.Command("emacs", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(f.Name())
}
