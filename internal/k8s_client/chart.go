package k8s_client

type Chart struct {
	Name       string
	Repository string
	Namespace  string
}

func InitChart(name, repo, namespace string) *Chart {
	chart := &Chart{
		Name:       name,
		Repository: repo,
		Namespace:  namespace,
	}
	return chart
}

func (c *Chart) AddRepo() error {
	return nil
}

func (c *Chart) Validate() error {
	return nil
}

func (c *Chart) Deploy() error {
	return nil
}
