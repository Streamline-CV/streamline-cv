package git

import (
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/rs/zerolog/log"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func getFileContentFromWorkingTree(repoPath, filePath string) ([]byte, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	fileContent, err := os.ReadFile(w.Filesystem.Join(filePath))
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func getFileContentFromBranch(repoPath, filePath, branchName string) ([]byte, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	file, err := tree.File(filePath)
	if err != nil {
		return nil, err
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Error().Msgf("Failed to close reader %e", err)
		}
	}(reader)

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetFileChange(filename string) (*api.FileState, error) {
	main, err := getFileContentFromBranch(".", filename, "main")
	if err != nil {
		return nil, err
	}
	current, err := getFileContentFromWorkingTree(".", filename)
	if err != nil {
		return nil, err
	}
	return &api.FileState{
		Source: main,
		Target: current,
	}, nil
}
