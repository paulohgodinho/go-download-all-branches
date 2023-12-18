package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	fmt.Println("Starting")

	const outputFolder = "./output"
	const mainRepoFolder = outputFolder + "/mainRepository"
	const repoURL = "https://github.com/go-git/go-git"

	os.RemoveAll(outputFolder)

	main, err := git.PlainClone(mainRepoFolder, false, &git.CloneOptions{
		URL:        repoURL,
		Progress:   os.Stdout,
		NoCheckout: false,
	})
	CheckIfError(err)
	fmt.Println("Clone Done")
	main.Fetch(&git.FetchOptions{
		Progress: os.Stdout,
	})
	fmt.Println("Fetch Done")

	allRemoteBranches, err := GetAllRemoteBranches(main)
	CheckIfError(err)

	for _, branchRef := range allRemoteBranches {
		fmt.Println(branchRef)

		shortBranchName := branchRef.Name().String()[len("refs/heads/"):]
		targetPath := "output/" + shortBranchName
		_, err := git.PlainClone(targetPath, false, &git.CloneOptions{
			URL:          repoURL,
			SingleBranch: true,

			ReferenceName: branchRef.Name(),
		})
		CheckIfError(err)
	}
}

func SetRepoToBranch(repo *git.Repository, targetBranch string) error {
	repo.Fetch(&git.FetchOptions{
		Progress: os.Stdout,
	})

	ref, err := GetRefForRemoteBranch(repo, "refs/heads/"+targetBranch)
	if err != nil {
		return fmt.Errorf("ref not found")
	}

	workTree, _ := repo.Worktree()
	err = workTree.Checkout(&git.CheckoutOptions{
		Hash:   ref.Hash(),
		Branch: plumbing.NewBranchReferenceName(targetBranch),
		Create: true,
	})
	CheckIfError(err)

	return nil
}

func GetRefForRemoteBranch(repo *git.Repository, refName string) (*plumbing.Reference, error) {
	allRefs, err := GetAllRemotesRefs(repo)
	CheckIfError(err)

	for _, val := range allRefs {
		if val.Name().String() == refName {
			return val, nil
		}
	}

	return nil, fmt.Errorf("ref not found")
}

func GetAllRemoteBranches(repo *git.Repository) ([]*plumbing.Reference, error) {
	allRefs, err := GetAllRemotesRefs(repo)
	CheckIfError(err)

	const refPrefix = "refs/heads/"
	var allRemoteBranches []*plumbing.Reference
	for _, ref := range allRefs {
		refName := ref.Name().String()
		if strings.HasPrefix(refName, refPrefix) == false {
			continue
		}
		allRemoteBranches = append(allRemoteBranches, ref)
	}

	return allRemoteBranches, nil
}

func GetAllRemotesRefs(repo *git.Repository) ([]*plumbing.Reference, error) {

	remote, err := repo.Remote("origin")
	if err != nil {
		panic(err)
	}
	refList, err := remote.List(&git.ListOptions{})
	if err != nil {
		panic(err)
	}

	return refList, nil
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
