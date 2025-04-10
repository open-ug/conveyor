package buildpacksdriver

import (
	"context"
	"fmt"
	"log"

	craneTypes "github.com/open-ug/conveyor/pkg/types"
	"github.com/buildpacks/pack/pkg/client"
)

// CreateBuildpacksImage creates a buildpacks image
func CreateBuildpacksImage(app *craneTypes.Application) error {
	// clone Source code
	fmt.Println("cloning repo")
	err := CloneGitRepo(app)

	if err != nil {
		log.Fatalf("failed to create pack client: %v", err)
		return err
	}
	fmt.Println("building image")

	err = BuildImage(app)

	if err != nil {
		log.Fatalf("failed to create pack client: %v", err)
		return err
	}

	return nil
}

func BuildImage(app *craneTypes.Application) error {
	cli, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to create pack client: %v", err)
		return err
	}

	buildOpts := client.BuildOptions{
		AppPath: "/usr/local/crane/git/" + app.Name,
		Builder: "heroku/builder:24",
		Image:   app.Name + "-bpimage",
		//PullPolicy: config.PullAlways,
	}

	if err := cli.Build(context.Background(), buildOpts); err != nil {
		log.Fatalf("failed to build image: %v", err)
		return err
	}

	return nil
}
