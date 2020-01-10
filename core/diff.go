package core

import (
	"fmt"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Diff(repo_path, from, to string, logger *Logger) ([]string, error) {
	files := make(map[string]bool)
	repo_path, err := filepath.Abs(repo_path)
	r, err := git.PlainOpen(repo_path)
	if err != nil {
		panic(err)
	}

	// ... retrieves the commit history
	commits, err := r.Log(&git.LogOptions{From: plumbing.NewHash(from)})

	if err != nil {
		return generateKeys(files), err
	}
	defer commits.Close()

	worktree, _ := r.Worktree()

	patterns, err := gitignore.ReadPatterns(worktree.Filesystem, nil)
	gitignoreMatcher := gitignore.NewMatcher(patterns)
	var prevCommit *object.Commit
	var prevTree *object.Tree

	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}

		fmt.Fprintf(logger, "Commit Msg: %s \n", commit.String())
		currentTree, err := commit.Tree()
		if err != nil {
			continue
		}

		if prevCommit == nil {
			prevCommit = commit
			prevTree = currentTree
			continue
		}

		changes, err := currentTree.Diff(prevTree)

		if err != nil {
			return generateKeys(files), err
		}

		for _, c := range changes {

			fmt.Fprintf(logger, "Changed File %v \n", c.To.Name)
			path := c.To.Name

			if gitignoreMatcher.Match(filepath.SplitList(path), false) {
				fmt.Fprintf(logger, "ignored file '%s' since it was in .gitignore", path)
				continue
			}
			if _, exists := files[path]; !exists {
				files[path] = true
			}
		}

		if commit.Hash.String() == plumbing.NewHash(to).String() {
			break
		}

		prevCommit = commit
		prevTree = currentTree
	}

	return generateKeys(files), nil
}

func generateKeys(m map[string]bool) []string {

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func recoverGettingChangedFiles(c *object.Change, logger *Logger) {
	if r := recover(); r != nil {
		fmt.Fprintf(logger, "Failed to fetch files from %s \n", c.From.Tree.Hash)
	}
}
