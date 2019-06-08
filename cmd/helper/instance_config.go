package helper

import (
	"math/rand"
	"time"
)

type InstanceConfig int

const (
	InstanceConfigUndefined InstanceConfig = iota
	InstanceConfigEur3
	InstanceConfigNamEurAsia1
	InstanceConfigNam3
	InstanceConfigNam6
	InstanceConfigRegionalAsiaEast1
	InstanceConfigRegionalAsiaEast2
	InstanceConfigRegionalAsiaNortheast1
	InstanceConfigRegionalAsiaNortheast2
	InstanceConfigRegionalAsiaSouth1
	InstanceConfigRegionalAsiaSoutheast1
	InstanceConfigRegionalAustraliaSoutheast1
	InstanceConfigRegionalEuropeNorth1
	InstanceConfigRegionalEuropeWest1
	InstanceConfigRegionalEuropeWest2
	InstanceConfigRegionalEuropeWest4
	InstanceConfigRegionalEuropeWest6
	InstanceConfigRegionalNorthamericaNortheast1
	InstanceConfigRegionalUsCentral1
	InstanceConfigRegionalUsEast1
	InstanceConfigRegionalUsEast4
	InstanceConfigRegionalUsWast1
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

func FindInstanceConfigByName(name string) InstanceConfig {
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
