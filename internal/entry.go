package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v31"
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

func (app *App) Create(dir string) error {
	buf := make([]byte, 0)
	conbuf := bytes.NewBuffer(buf)

	now := time.Now()

	fname := fmt.Sprintf("%s.md", now.Format("2006-01-02"))
	if dir != "" {
		dirPath := fmt.Sprintf("%s/%s", app.GitDir, dir)
		_ = os.Mkdir(dirPath, 0700)
		fname = fmt.Sprintf("%s/%s", dir, fname)
	}
	fpath := fmt.Sprintf("%s/%s", app.GitDir, fname)

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
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

	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	conbuf.Write(endOfContent)
	if _, err := f.Write(conbuf.Bytes()); err != nil {
		return err
	}

	if err := app.commit(fname, fmt.Sprintf("created file %s", fname)); err != nil {
		return err
	}

	return nil
}

func (app *App) Update(date string) error {

	if !strings.ContainsAny(date, ".md") { //support not adding extension
		date = date + ".md"
	}

	fname := date

	tf, err := ioutil.TempFile("/tmp", "entry*.md")
	if err != nil {
		return err
	}
	defer tf.Close()
	defer os.Remove(tf.Name())

	fpath := fmt.Sprintf("%s/%s", app.GitDir, date)
	exFile, err := os.Open(fpath)
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

	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) //OS truncate
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return err
	}

	if err := app.commit(fname, fmt.Sprintf("updated file %s", fpath)); err != nil {
		return err
	}

	return nil
}

func (app *App) List() error {
	return filepath.Walk(app.GitDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dir, fname := filepath.Split(path)
			if strings.Contains(".git", dir) {
				return nil
			}
			dir = strings.Replace(dir, app.GitDir+"/", "", 1)
			if fname != "README.md" && filepath.Ext(fname) == ".md" {
				name := strings.ReplaceAll(fname, ".md", "")
				if dir != "" {
					name = fmt.Sprintf("%s%s", dir, name)
				}
				fmt.Println(name)
			}
		}
		return nil
	})
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

func (app *App) commit(path, msg string) error {
	var signature = &git.Signature{ //TODO config
		Name:  "Joel Holmes",
		Email: "holmes89@gmail.com",
		When:  time.Now(),
	}

	repo, err := git.OpenRepository(app.GitDir)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	idx, err := repo.Index()
	if err != nil {
		return err
	}

	if err := idx.AddByPath(path); err != nil {
		return err
	}

	treeID, err := idx.WriteTree()
	if err != nil {
		return err
	}

	if err := idx.Write(); err != nil {
		return err
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return err
	}

	commitTarget, err := repo.LookupCommit(head.Target())
	if err != nil {
		return err
	}

	if _, err := repo.CreateCommit("refs/heads/master", signature, signature, msg, tree, commitTarget); err != nil {
		return err
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	if err := remote.Push([]string{"refs/heads/master"}, &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      app.credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}); err != nil { // working off of master for now, maybe move to branches in the future?
		return err
	}

	return nil
}
