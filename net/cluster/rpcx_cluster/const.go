package cherryRpcxCluster

import (
	"fmt"
)

const (
	serviceNameFormat = "nodes.%s.%s"
)

func getServiceName(nodeType string, nodeId string) string {
	return fmt.Sprintf(serviceNameFormat, nodeType, nodeId)
}
