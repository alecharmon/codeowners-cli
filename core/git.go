package core

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Diff(repo_path, from, to string) ([]string, error) {
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

		fmt.Println(commit.Hash.String())
		currentTree, err := commit.Tree()
		if err != nil {
			continue
			return generateKeys(files), err
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
			for _, path := range getChangedFiles(c) {
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

func getChangedFiles(c *object.Change) (res []string) {
	defer recoverGettingChangedFiles(c)
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

func recoverGettingChangedFiles(c *object.Change) {
	if r := recover(); r != nil {
		fmt.Println("Failed to fetch files from ", c.From.Name)
	}
}
