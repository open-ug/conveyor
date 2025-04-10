package buildpacksdriver

import (
	"fmt"

	craneTypes "github.com/open-ug/conveyor/pkg/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// CloneGitRepo clones a git repository
func CloneGitRepo(app *craneTypes.Application) error {
	_, err := git.PlainClone("/usr/local/crane/git/"+app.Name, false, &git.CloneOptions{
		URL: app.Spec.Source.GitRepo.URL,
		Auth: &http.BasicAuth{
			Username: app.Spec.Source.GitRepo.Username,
			Password: app.Spec.Source.GitRepo.Password,
		},
	})
	if err != nil {
		return fmt.Errorf("error cloning git repository: %v", err)
	}
	fmt.Println("Cloned")

	return nil
}

// CreateBuildpacksImage creates a buildpacks image
