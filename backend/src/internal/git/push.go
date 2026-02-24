package git

import (
	"fmt"
	"time"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func PushRepo(token, repoName, localPath string) error {
	// 1️⃣ Init git repo
	repo, err := git.PlainInit(localPath, false)
	if err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}
	owner, err := GetAuthenticatedUser(token)
	if err != nil {
		return err
	}
	// 2️⃣ Add remote origin
	remoteURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repoName)

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	})
	if err != nil {
		return fmt.Errorf("add remote failed: %w", err)
	}

	// 3️⃣ Stage all files
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// 4️⃣ Commit
	_, err = worktree.Commit("Initial commit from platform", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Platform Bot",
			Email: "platform@company.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	// 5️⃣ Push
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "x-access-token", // required by GitHub
			Password: token,
		},
	})

	if err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}




