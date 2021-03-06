package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

func Diff(repo_path, from, to string, newFilesOnly bool, logger *Logger) ([]string, []error) {
	files := make(map[string]bool)
	repo_path, err := filepath.Abs(repo_path)
	r, err := git.PlainOpen(repo_path)
	if err != nil {
		return nil, []error{err}
	}

	worktree, _ := r.Worktree()
	patterns, err := gitignore.ReadPatterns(worktree.Filesystem, nil)
	gitignoreMatcher := gitignore.NewMatcher(patterns)

	if to == "master" && from == "head" {
		return getAllFiles(worktree.Filesystem.Root(), gitignoreMatcher, logger)
	}

	// ... retrieves the commit history
	commits, err := r.Log(&git.LogOptions{From: plumbing.NewHash(from)})

	if err != nil {
		return generateKeys(files), []error{err}
	}
	defer commits.Close()

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
			return generateKeys(files), []error{err}
		}

		for _, c := range changes {
			action, err := c.Action()
			if err != nil {
				continue
			}
			if newFilesOnly && (action != merkletrie.Insert) {
				continue
			}
			if action == merkletrie.Delete {
				continue
			}
			fmt.Fprintf(logger, "Changed File %v \n", c.To.Name)
			path := c.To.Name
			if len(path) < 1 {
				continue
			}
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

func getAllFiles(path string, matcher gitignore.Matcher, logger *Logger) ([]string, []error) {
	files := make(map[string]bool)

	fmt.Fprintf(logger, "Getting all files")
	pathL := len(path) + 1
	err := filepath.Walk(path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			p = string([]rune(p)[pathL:])

			if strings.Contains(p, ".git/") || strings.Contains(p, ".git/") {
				return nil
			}
			if matcher.Match(filepath.SplitList(p), false) {
				fmt.Fprintf(logger, "ignored file '%s' since it was in .gitignore", path)
				return nil
			}

			if _, exists := files[p]; !exists {
				files[p] = true
			}
			return nil
		})
	return generateKeys(files), []error{err}
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
