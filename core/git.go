package core

import (
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
			c.To.Tree.Files().ForEach(func(f *object.File) error {
				if _, exists := files[f.Name]; !exists {
					files[f.Name] = true
				}
				return nil
			})
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
