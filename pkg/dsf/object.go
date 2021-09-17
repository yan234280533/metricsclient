package dsf

type DataSourceObjectType string

const (
	DataSourceObjectNode      DataSourceObjectType = "node"
	DataSourceObjectPod       DataSourceObjectType = "pod"
	DataSourceObjectContainer DataSourceObjectType = "container"
)

type DataSourceObjectName struct {
	DataSourceObjectType
	NodeName      string
	ContainerName string
	PodName       string
	Namespace     string
}

func NewNodeDataSourceObject(nodeName string) DataSourceObjectName {
	return DataSourceObjectName{DataSourceObjectType: DataSourceObjectNode, NodeName: nodeName}
}

func NewPodDataSourceObject(podName string, namespace string) DataSourceObjectName {
	return DataSourceObjectName{DataSourceObjectType: DataSourceObjectPod, Namespace: namespace, PodName: podName}
}

func NewContainerDataSourceObject(podName string, namespace string, containerName string) DataSourceObjectName {
	return DataSourceObjectName{DataSourceObjectType: DataSourceObjectContainer, Namespace: namespace, PodName: podName, ContainerName: containerName}
}

func IsNodeDataSourceObject(name DataSourceObjectName) bool {
	return name.DataSourceObjectType == DataSourceObjectNode
}

func IsPodDataSourceObject(name DataSourceObjectName) bool {
	return name.DataSourceObjectType == DataSourceObjectPod
}

func IsContainerDataSourceObject(name DataSourceObjectName) bool {
	return name.DataSourceObjectType == DataSourceObjectContainer
}
