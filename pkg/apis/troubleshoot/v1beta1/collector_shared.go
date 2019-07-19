package v1beta1

type ClusterInfo struct {
}

type ClusterResources struct {
}

type Secret struct {
	Name         string `json:"name" yaml:"name"`
	Namespace    string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Key          string `json:"key,omitempty" yaml:"key,omitempty"`
	IncludeValue bool   `json:"includeValue,omitempty" yaml:"includeValue,omitempty"`
}

type LogLimits struct {
	MaxAge   string `json:"maxAge,omitempty" yaml:"maxAge,omitempty"`
	MaxLines int64  `json:"maxLines,omitempty" yaml:"maxLines,omitempty"`
}

type Logs struct {
	Selector  []string   `json:"selector" yaml:"selector"`
	Namespace string     `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Limits    *LogLimits `json:"limits,omitempty" yaml:"omitempty"`
}

type Run struct {
	Name            string   `json:"name" yaml:"name"`
	Namespace       string   `json:"namespace" yaml:"namespace"`
	Image           string   `json:"image" yaml:"image"`
	Command         []string `json:"command,omitempty" yaml:"command,omitempty"`
	Args            []string `json:"args,omitempty" yaml:"args,omitempty"`
	Timeout         string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	ImagePullPolicy string   `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
}

type Copy struct {
	Selector      []string `json:"selector" yaml:"selector"`
	Namespace     string   `json:"namespace" yaml:"namespace"`
	ContainerPath string   `json:"containerPath" yaml:"containerPath"`
	ContainerName string   `json:"containerName,omitempty" yaml:"containerName,omitempty"`
}

type Collect struct {
	ClusterInfo      *ClusterInfo      `json:"clusterInfo,omitempty" yaml:"clusterInfo,omitempty"`
	ClusterResources *ClusterResources `json:"clusterResources,omitempty" yaml:"clusterResources,omitempty"`
	Secret           *Secret           `json:"secret,omitempty" yaml:"secret,omitempty"`
	Logs             *Logs             `json:"logs,omitempty" yaml:"logs,omitempty"`
	Run              *Run              `json:"run,omitempty" yaml:"run,omitempty"`
	Copy             *Copy             `json:"copy,omitempty" yaml:"copy,omitempty"`
}
