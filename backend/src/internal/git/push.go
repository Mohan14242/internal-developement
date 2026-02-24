package git

import (
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func PushRepo(token, repoName, localPath, branch string) error {
	// 1️⃣ Init repo
	repo, err := git.PlainInit(localPath, false)
	if err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}

	owner, err := GetAuthenticatedUser(token)
	if err != nil {
		return err
	}

	// 2️⃣ Worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// 3️⃣ Add files
	if _, err := worktree.Add("."); err != nil {
		return err
	}

	// 4️⃣ FIRST COMMIT (on default branch)
	_, err = worktree.Commit("Initial commit from platform", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Platform Bot",
			Email: "platform@company.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// 5️⃣ Create & checkout target branch (dev)
	refName := plumbing.NewBranchReferenceName(branch)

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: refName,
		Create: true,
	})
	if err != nil {
		return fmt.Errorf("create/checkout branch %s failed: %w", branch, err)
	}

	// 6️⃣ Add remote
	remoteURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repoName)

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	})
	if err != nil {
		return err
	}

	// 7️⃣ Push dev branch
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(
				fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch),
			),
		},
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
	})
	if err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}