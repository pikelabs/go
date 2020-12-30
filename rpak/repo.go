package rpak

import (
	"errors"

	git "github.com/go-git/go-git"
)

type Repo interface {
	Open(path string) error
	Create(path string) error
	Get(url, path string) error
	IsClean() (bool, error)
	Stage(files ...string) error
}

type GitRepo struct {
	repo *git.Repository
}

func (r *GitRepo) Open(path string) (err error) {
	r.repo, err = git.PlainOpen(path)
	return err
}

func (r *GitRepo) Create(path string) (err error) {
	r.repo, err = git.PlainInit(path, false)
	return err
}

func (r *GitRepo) Get(url, path string) error {
	return nil
}

func (r *GitRepo) IsClean() (bool, error) {
	if r == nil || r.repo == nil {
		return false, errors.New("not initialized")
	}

	wt, err := r.repo.Worktree()
	if err != nil {
		return false, err
	}
	s, err := wt.Status()
	if err != nil {
		return false, err
	}
	return s.IsClean(), nil
}

func (r *GitRepo) Stage(paths ...string) error {
	if r == nil || r.repo == nil {
		return errors.New("not initialized")
	}

	wt, err := r.repo.Worktree()
	if err != nil {
		return err
	}
	for _, path := range paths {
		_, err = wt.Add(path)
		if err != nil {
			return err
		}
	}
	return nil
}
