package digitalocean

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	logger "github.com/sirupsen/logrus"
	"time"
)

func convertLabelsIntoTags(labels map[string]string) []string {
	var tags []string
	if len(labels) > 0 {
		for k, v := range labels {
			composed := fmt.Sprintf("%%s", k, v)
			tags = append(tags, composed)
		}
	}
	return tags
}
func (g *DigitalOceanProvider) Create(clusterp kinds.ProviderCluster) error {

	//Convert generic node pool type into a specific DO resource
	var convertedNodePool []*godo.KubernetesNodePoolCreateRequest

	for _, model := range clusterp.NodePools {
		nodePool := new(godo.KubernetesNodePoolCreateRequest)
		nodePool.Name = model.NodePool.Name
		nodePool.Count = model.NodePool.Count
		nodePool.Size = model.NodePool.NodeType
		//	nodePool.Tags = convertLabelsIntoTags(model.NodePool.Labels)
		convertedNodePool = append(convertedNodePool, nodePool)
	}

	req := &godo.KubernetesClusterCreateRequest{
		Name:        clusterp.Name,
		RegionSlug:  clusterp.Region,
		VersionSlug: clusterp.KubernetesVersion,
		NodePools:   convertedNodePool,
		//		Tags:        convertLabelsIntoTags(clusterp.Labels),
	}

	_, resp, err := g.ClusterManagerClient.Kubernetes.Create(g.Context,
		req)
	if err != nil {
		logger.Error(resp)
		return err
	}

	if resp.StatusCode == 201 {
		for {
			clust, _, err :=
				g.ClusterManagerClient.Kubernetes.Get(g.Context, clusterp.Name)
			if err != nil {
				continue
			}
			if clust.Status.State == "running" {
				color.Green("Cluster running")
				return nil
			}
			time.Sleep(time.Second)
		}
	}
	return nil
}
