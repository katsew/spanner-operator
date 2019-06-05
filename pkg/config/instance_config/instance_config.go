package instance_config

import (
	"math/rand"
	"time"
)

type InstanceConfig int

const (
	Undefined InstanceConfig = iota
	Eur3
	NamEurAsia1
	Nam3
	Nam6
	RegionalAsiaEast1
	RegionalAsiaEast2
	RegionalAsiaNortheast1
	RegionalAsiaNortheast2
	RegionalAsiaSouth1
	RegionalAsiaSoutheast1
	RegionalAustraliaSoutheast1
	RegionalEuropeNorth1
	RegionalEuropeWest1
	RegionalEuropeWest2
	RegionalEuropeWest4
	RegionalEuropeWest6
	RegionalNorthamericaNortheast1
	RegionalUsCentral1
	RegionalUsEast1
	RegionalUsEast4
	RegionalUsWast1
)

var instanceConfigs = [22]string{
	"undefined",
	"eur3",
	"nam-eur-asia1",
	"nam3",
	"nam6",
	"regional-asia-east1",
	"regional-asia-east2",
	"regional-asia-northeast1",
	"regional-asia-northeast2",
	"regional-asia-south1",
	"regional-asia-southeast1",
	"regional-australia-southeast1",
	"regional-europe-north1",
	"regional-europe-west1",
	"regional-europe-west2",
	"regional-europe-west4",
	"regional-europe-west6",
	"regional-northamerica-northeast1",
	"regional-us-central1",
	"regional-us-east1",
	"regional-us-east4",
	"regional-us-west1",
}

func (s InstanceConfig) String() string {
	return instanceConfigs[s]
}

func FindByName(name string) InstanceConfig {
	for i, c := range instanceConfigs {
		if c == name {
			return InstanceConfig(i)
		}
	}
	return InstanceConfig(0)
}

func GetRandomInstanceConfig() string {
	min := 1
	max := len(instanceConfigs) - 1
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn((max - min) + min)
	return instanceConfigs[i]
}
