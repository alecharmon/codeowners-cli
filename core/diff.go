package core

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Diff(repo_path, from, to string, logger *Logger) ([]string, error) {
	files := make(map[string]bool)
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

	var prevCommit *object.Commit
	var prevTree *object.Tree

	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}

		fmt.Fprintf(logger, "HASH: %s \n", commit.Hash.String())
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
			for _, path := range getChangedFiles(c, logger) {
				if _, exists := files[path]; !exists {
					files[path] = true
				}
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

func getChangedFiles(c *object.Change, logger *Logger) (res []string) {
	defer recoverGettingChangedFiles(c, logger)
	c.To.Tree.Files().ForEach(func(f *object.File) error {
		res = append(res, f.Name)
		return nil
	})
	return res
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
