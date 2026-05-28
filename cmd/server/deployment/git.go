package deployment

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitFetchStrategy interface {
	Fetch(repoURL, destDir string) error
}

type BranchStrategy struct {
	Branch string
}

func (s *BranchStrategy) Fetch(repoURL, destDir string) error {
	_, err := gogit.PlainClone(destDir, false, &gogit.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(s.Branch),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRepoCloneFailed, err)
	}
	return nil
}

type TagStrategy struct {
	Tag string
}

func (s *TagStrategy) Fetch(repoURL, destDir string) error {
	_, err := gogit.PlainClone(destDir, false, &gogit.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewTagReferenceName(s.Tag),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRepoCloneFailed, err)
	}
	return nil
}

type CommitStrategy struct {
	Commit string
}

func (s *CommitStrategy) Fetch(repoURL, destDir string) error {
	repo, err := gogit.PlainClone(destDir, false, &gogit.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRepoCloneFailed, err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCommitCheckoutFailed, err)
	}

	err = w.Checkout(&gogit.CheckoutOptions{
		Hash: plumbing.NewHash(s.Commit),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCommitCheckoutFailed, err)
	}
	return nil
}
