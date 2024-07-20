package models

type Deployment struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Namespace   string `yaml:"namespace"`
		Annotations struct {
			FluxcdIoAutomated                  string `yaml:"fluxcd.io/automated"`
			FluxWeaveWorksTagMsFreeresourceDiy string `yaml:"flux.weave.works/tag.service-name"`
		} `yaml:"annotations,omitempty"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas,omitempty"`
		Selector struct {
			MatchLabels struct {
				App     string `yaml:"app"`
				Version string `yaml:"version"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels struct {
					App     string `yaml:"app"`
					Version string `yaml:"version"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				ImagePullSecrets []struct {
					Name string `yaml:"name"`
				} `yaml:"imagePullSecrets"`
				Containers []struct {
					Ports []struct {
						ContainerPort int `yaml:"containerPort,omitempty"`
					} `yaml:"ports,omitempty"`
					Name            string    `yaml:"name"`
					Image           string    `yaml:"image"`
					Resources       Resources `yaml:"resources"`
					Env             Env       `yaml:"env"`
					EnvFrom         EnvFrom   `yaml:"envFrom"`
					ImagePullPolicy string    `yaml:"imagePullPolicy"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type Service struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Labels    struct {
			App string `yaml:"app"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Ports []struct {
			Name       string `yaml:"name"`
			Protocol   string `yaml:"protocol"`
			Port       int    `yaml:"port"`
			TargetPort int    `yaml:"targetPort"`
			NodePort   int    `yaml:"nodePort,omitempty"`
		} `yaml:"ports"`
		Selector struct {
			App string `yaml:"app"`
		} `yaml:"selector"`
		Type string `yaml:"type"`
	} `yaml:"spec"`
}

type HorizontalPodAutoscaler struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		ScaleTargetRef struct {
			Kind       string `yaml:"kind"`
			Name       string `yaml:"name"`
			APIVersion string `yaml:"apiVersion"`
		} `yaml:"scaleTargetRef"`
		MinReplicas int `yaml:"minReplicas"`
		MaxReplicas int `yaml:"maxReplicas"`
		Metrics     []struct {
			Type     string `yaml:"type"`
			Resource struct {
				Name   string `yaml:"name"`
				Target struct {
					Type               string `yaml:"type"`
					AverageUtilization int    `yaml:"averageUtilization"`
				} `yaml:"target"`
			} `yaml:"resource"`
		} `yaml:"metrics"`
	} `yaml:"spec"`
}

type ImagePolicy struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		ImageRepositoryRef struct {
			Name string `yaml:"name"`
		} `yaml:"imageRepositoryRef"`
		FilterTags struct {
			Pattern string `yaml:"pattern"`
			Extract string `yaml:"extract"`
		} `yaml:"filterTags"`
		Policy struct {
			Numerical struct {
				Order string `yaml:"order"`
			} `yaml:"numerical"`
		} `yaml:"policy"`
	} `yaml:"spec"`
}

type ImageRepository struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Image         string `yaml:"image"`
		Interval      string `yaml:"interval"`
		CertSecretRef struct {
			Name string `yaml:"name"`
		} `yaml:"certSecretRef"`
		SecretRef struct {
			Name string `yaml:"name"`
		} `yaml:"secretRef"`
	} `yaml:"spec"`
}

type Payload struct {
	ServiceName      string    `json:"servicename"`
	Image            string    `json:"image"`
	Resources        Resources `json:"resources,omitempty"`
	Namespace        string    `json:"namespace"`
	Env              Env       `json:"environment"`
	ImagePullSecrets string    `json:"pullSecret"`
	Port             int       `json:"port,omitempty"`
	NodePort         int       `json:"nodePort,omitempty"`
	MinPods          int       `json:"minPods,omitempty"`
	MaxPods          int       `json:"maxPods,omitempty"`
	FluxTag          string    `json:"fluxTag,omitempty"`
}

type Resources struct {
	Limits   Resource `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests Resource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type Resource struct {
	Memory string `json:"memory,omitempty" yaml:"memory"`
	CPU    string `json:"cpu,omitempty" yaml:"cpu"`
}

type Env []struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type EnvFrom []struct {
	ConfigMapRef struct {
		Name string `json:"name" yaml:"name"`
	} `json:"configMapRef" yaml:"configMapRef"`
}
