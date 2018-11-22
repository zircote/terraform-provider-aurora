package aurora

import (
	"github.com/paypal/gorealis"
	"log"
)

type Config struct {
	Username     string
	Password     string
	Transport    string
	Debug        bool
	CaCertsPath  string
	ClientKey    string
	ClientCert   string
	ZkUrl        string
	SchedulerUrl string
	Timeout      int

	client  realis.Realis
	cluster realis.Cluster
}

func (c *Config) CreateAuroraClient() error {
	clientOptions := []realis.ClientOption{}
	if len(c.Username) > 0 {
		clientOptions = append(clientOptions, realis.BasicAuth(c.Username, c.Password))
		log.Printf("Aurora Provider authenticating with Aurora [Username: %s, Password: %s",
			c.Username, c.Password)
	}

	if c.ZkUrl != "" {
		clientOptions = append(clientOptions, realis.ZKUrl(c.ZkUrl))
	} else {
		clientOptions = append(clientOptions, realis.SchedulerUrl(c.SchedulerUrl))
	}

	if c.CaCertsPath != "" {
		clientOptions = append(clientOptions, realis.Certspath(c.CaCertsPath))
	}

	if c.ClientKey != "" && c.ClientCert != "" {
		clientOptions = append(clientOptions, realis.ClientCerts(c.ClientKey, c.ClientCert))
	}

	if c.cluster.ZK != "" {
		clientOptions = append(clientOptions, realis.ZKCluster(&c.cluster))
	}

	if c.Timeout > 0 {
		clientOptions = append(clientOptions, realis.TimeoutMS(c.Timeout))
	} else {
		clientOptions = append(clientOptions, realis.TimeoutMS(20000))
	}
	if c.Debug {
		clientOptions = append(clientOptions, realis.SetLogger(realis.LevelLogger{}))
	} else {
		clientOptions = append(clientOptions, realis.SetLogger(realis.NoopLogger{}))
	}

	gr, err := realis.NewRealisClient(clientOptions...)
	if err != nil {
		log.Print(err)
		return err
	}
	c.client = gr
	return nil
}
