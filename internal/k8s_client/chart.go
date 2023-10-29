package k8s_client

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	helm "github.com/mittwald/go-helm-client"
)

type Chart struct {
	HelmClient helm.Client
	Spec       helm.ChartSpec
	Output     *bytes.Buffer
}

func InitChart(name, namespace string) (*Chart, error) {

	opt := &helm.Options{
		Namespace:        namespace, // Change this to the namespace you wish the client to operate in.
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         func(format string, v ...interface{}) {},
	}

	chartsRepoURL := os.Getenv("CHARTS_REPO_URL")
	chartsRepoPath := os.Getenv("CHARTS_REPO_PATH")

	spec := helm.ChartSpec{}

	if name != "" {
		spec = helm.ChartSpec{
			Namespace:       namespace,
			ReleaseName:     namespace,
			CreateNamespace: true,
			DryRun:          false,
			ChartName:       fmt.Sprintf("%s/%s/%s.tgz", chartsRepoURL, chartsRepoPath, name),
		}
	}

	helmClient, err := helm.New(opt)
	if err != nil {
		panic(err)
	}
	chart := &Chart{
		HelmClient: helmClient,
		Spec:       spec,
	}
	return chart, nil
}

func (c *Chart) AddRepo() error {
	return nil
}

func (c *Chart) Validate() error {
	return nil
}

func (c *Chart) Deploy(values map[string]string) error {
	for key, value := range values {
		c.Spec.ValuesOptions.Values = append(
			c.Spec.ValuesOptions.Values,
			fmt.Sprintf("%s=%s", key, value),
		)
	}
	release, err := c.HelmClient.InstallOrUpgradeChart(context.Background(), &c.Spec, nil)
	log.Println(release)
	return err
}

func (c *Chart) Destroy(release string) error {
	return c.HelmClient.UninstallReleaseByName(release)
}
