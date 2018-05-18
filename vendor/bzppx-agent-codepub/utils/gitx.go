package utils

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bytes"
	"log"
	"os/exec"
	"runtime"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var branchNamePrefix = "refs/heads/auto-"

type GitX struct{}

func NewGitX() *GitX {
	return &GitX{}
}

type GitXParams struct {
	Url        string `json:"url"`
	SshKey     string `json:"ssh_key"`
	SshKeySalt string `json:"ssh_key_salt"`
	Path       string `json:"path"`
	Branch     string `json:"branch"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DirUser    string `json:"dir_user"`
}

// 验证参数合法性
func (g *GitX) Validate(params GitXParams) error {
	if params.Path == "" {
		return errors.New("Gitx params error: path requied!")
	}
	if params.Url == "" {
		return errors.New("Gitx params error: url requied!")
	}
	if g.IsHTTP(params) {
		if (params.Username != "" && params.Password == "") ||
			(params.Username == "" && params.Password != "") {
			return errors.New("Gitx params error: username and password requied!")
		}
	} else if params.SshKey == "" {
		return errors.New("Gitx params error: ssh_key requied!")
	}
	return nil
}

// 切换分支
func (g *GitX) Checkout(name string, params GitXParams) (err error) {
	r, err := git.PlainOpen(params.Path)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + name),
		Force:  true,
	})
	return
}

// 清除分支
func (g *GitX) CleanBranch(params GitXParams) (err error) {
	//var hash string
	//var hashObj plumbing.Hash
	var r *git.Repository
	//var ref *plumbing.Reference
	if r, err = git.PlainOpen(params.Path); err != nil {
		return
	}
	refs, err := r.References()
	if err != nil {
		return
	}
	h, err := r.Head()
	if err != nil {
		return
	}
	headBranchName := h.Name().String()
	refs.ForEach(func(ref0 *plumbing.Reference) (e error) {
		//fmt.Println(ref0.Name().String(), ref0.Hash().String())
		if headBranchName != ref0.Name().String() && strings.HasPrefix(ref0.Name().String(), "refs/heads/") {
			err = r.Storer.RemoveReference(ref0.Name())
			if err != nil {
				e = err
				return
			}
		}
		return nil
	})
	return
}

// 创建分支名称
func (g *GitX) CreateBranchName(params GitXParams) (name string, err error) {
	var hash string
	// var hashObj plumbing.Hash
	var ref *plumbing.Reference
	hash, _, _, ref, err = g.GetHash(params)
	if err != nil {
		return
	}
	f := ""
	if ref != nil {
		f = filepath.Base(ref.Name().Short())
	} else {
		f = hash
	}
	name = branchNamePrefix + f + "-" + NewStr().NewLenChars(8, StdChars)
	return
}

// 发布代码
func (g *GitX) Publish(params GitXParams) (commitId string, err error) {
	defer func() {
		go func() {
			if runtime.GOOS != "windows" {
				cmd := exec.Command("chown", "-R", params.DirUser, params.Path)
				var out bytes.Buffer
				cmd.Stderr = &out
				err = cmd.Run()
				if err != nil {
					log.Println(out.String())
				}
			}
		}()
	}()

	if NewFile().PathIsEmpty(params.Path) {
		err = nil
		_, err = g.Clone(params)
	} else {
		p := filepath.Join(params.Path, ".git")
		if _, err := os.Stat(p); err != nil || NewFile().PathIsEmpty(p) {
			return "", errors.New("Gitx publish error: target path is not a git repository!")
		}
		_, err = git.PlainOpen(params.Path)
		if err == nil {
			_, err = g.Fetch(params)
		}
	}
	if err != nil {
		return
	}
	branchShortName := ""
	branchShortName, _, err = g.CreateBranch(params)
	if err != nil {
		return
	}
	//output := ""
	err = g.Checkout(branchShortName, params)
	if err != nil {
		return
	}
	err = g.CleanBranch(params)
	if err != nil {
		return
	}
	// get commit_id
	commitId, err = g.LastCommitId(params)

	return
}

// 创建分支
func (g *GitX) CreateBranch(params GitXParams) (branchShortName, branchName string, err error) {
	_, hashObj, r, _, err := g.GetHash(params)
	if err != nil {
		return
	}
	branchName, err = g.CreateBranchName(params)
	if err != nil {
		return
	}
	branchShortName = filepath.Base(branchName)
	err = r.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName(branchName), hashObj))
	return
}

// 获取 hash
func (g *GitX) GetHash(params GitXParams) (hash string, hashObj plumbing.Hash, r *git.Repository, ref *plumbing.Reference, err error) {
	if r, err = git.PlainOpen(params.Path); err != nil {
		return
	}
	refs, err := r.References()
	if err != nil {
		return
	}
	refs.ForEach(func(ref0 *plumbing.Reference) (e error) {
		//fmt.Println(ref0.Name().String(), ref0.Hash().String())
		if !strings.HasPrefix(ref0.Name().String(), "refs/heads/") {
			if params.Branch == filepath.Base(ref0.Name().String()) {
				hash = ref0.Hash().String()
				ref = ref0
				hashObj = ref0.Hash()
				return errors.New("")
			}
		}
		return nil
	})
	if hash == "" {
		var c *object.Commit
		c, err = r.CommitObject(plumbing.NewHash(params.Branch))
		if err == nil {
			hash = params.Branch
			hashObj = c.Hash
			return
		}
	}
	return
}

// 拉取代码
func (g *GitX) Fetch(params GitXParams) (r *git.Repository, err error) {
	if err = g.Validate(params); err != nil {
		return
	}
	if r, err = git.PlainOpen(params.Path); err != nil {
		return
	}
	rconfig, err := r.Storer.Config()
	if err != nil {
		return
	}
	remoteConfig := &config.RemoteConfig{
		Name:  git.DefaultRemoteName,
		URLs:  []string{params.Url},
		Fetch: []config.RefSpec{"+refs/heads/*:refs/remotes/" + git.DefaultRemoteName + "/*"},
	}
	rconfig.Remotes = map[string]*config.RemoteConfig{
		git.DefaultRemoteName: remoteConfig,
	}
	r.Storer.SetConfig(rconfig)
	var opt git.FetchOptions
	if opt, err = g.FetchOptions(params); err != nil {
		return
	}
	if err = r.Fetch(&opt); err == git.NoErrAlreadyUpToDate {
		err = nil
	}
	return
}

// 克隆代码
func (g *GitX) Clone(params GitXParams) (r *git.Repository, err error) {
	if err = g.Validate(params); err != nil {
		return
	}
	var opt git.CloneOptions
	opt, err = g.CloneOptions(params)
	if err != nil {
		return
	}
	r, err = git.PlainClone(params.Path, false, &opt)
	return
}

// 获取最近一次提交的 commit_id
func (g *GitX) LastCommitId(params GitXParams) (hash string, err error) {
	if err = g.Validate(params); err != nil {
		return
	}
	var r *git.Repository
	if r, err = git.PlainOpen(params.Path); err != nil {
		return
	}
	commitIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		return
	}
	commit, err := commitIter.Next()
	if err != nil {
		return
	}

	hash = commit.Hash.String()
	return
}

// 获取克隆参数
func (g *GitX) CloneOptions(params GitXParams) (opt git.CloneOptions, err error) {
	opt.URL = params.Url
	if g.IsNeedAuth(params) {
		opt.Auth, err = g.GetAuth(params)
		if err != nil {
			return
		}
	}
	return
}

// 获取拉取参数
func (g *GitX) FetchOptions(params GitXParams) (opt git.FetchOptions, err error) {
	opt.RemoteName = git.DefaultRemoteName
	opt.RefSpecs = []config.RefSpec{"+refs/heads/*:refs/remotes/" + git.DefaultRemoteName + "/*"}
	if g.IsNeedAuth(params) {
		opt.Auth, err = g.GetAuth(params)
		if err != nil {
			return
		}
	}
	return
}

// 获取 auth 信息
func (g *GitX) GetAuth(params GitXParams) (auth transport.AuthMethod, err error) {
	var signer ssh.Signer
	if g.IsHTTP(params) {
		customClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 30 * time.Minute,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		client.InstallProtocol("https", githttp.NewClient(customClient))
		client.InstallProtocol("http", githttp.NewClient(customClient))
		auth = &githttp.BasicAuth{
			Username: params.Username,
			Password: params.Password,
		}
	} else {
		if params.SshKeySalt != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(params.SshKey), []byte(params.SshKeySalt))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(params.SshKey))
		}
		if err != nil {
			return
		}
		auth = &gitssh.PublicKeys{
			User:   "git",
			Signer: signer,
			HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			},
		}

	}
	return
}

// 是否是 http
func (g *GitX) IsHTTP(params GitXParams) bool {
	return strings.HasPrefix(params.Url, "http://") || strings.HasPrefix(params.Url, "https://")
}

// 是否需要 auth
func (g *GitX) IsNeedAuth(params GitXParams) bool {
	if g.IsHTTP(params) {
		return params.Username != "" && params.Password != ""
	}
	return params.SshKey != ""
}
