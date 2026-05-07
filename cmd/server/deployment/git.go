package deployment

import (
	"fmt"
	"os"
	"os/exec"
)

type GitFetchStrategy interface {
	Fetch(repoURL, destDir string) error
}

type BranchStrategy struct {
	Branch string
}

func (s *BranchStrategy) Fetch(repoURL, destDir string) error {
	cmd := exec.Command("git", "clone", "-b", s.Branch, "--single-branch", repoURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type TagStrategy struct {
	Tag string
}

func (s *TagStrategy) Fetch(repoURL, destDir string) error {
	cmd := exec.Command("git", "clone", "-b", s.Tag, "--single-branch", repoURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type CommitStrategy struct {
	Commit string
}

func (s *CommitStrategy) Fetch(repoURL, destDir string) error {
	cloneCmd := exec.Command("git", "clone", repoURL, destDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrRepoCloneFailed, err)
	}

	checkoutCmd := exec.Command("git", "checkout", s.Commit)
	checkoutCmd.Dir = destDir
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr
	if err := checkoutCmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrCommitCheckoutFailed, err)
	}
	return nil
}
