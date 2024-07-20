package models

type Route struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Namespace   string `yaml:"namespace"`
		Annotations struct {
			DisableCookies string `yaml:"haproxy.router.openshift.io/disable_cookies"`
		} `yaml:"annotations,omitempty"`
	} `yaml:"metadata"`
	Spec struct {
		Host interface{} `yaml:"host"`
		Path string      `yaml:"path"`
		To   struct {
			Kind   string `yaml:"kind"`
			Name   string `yaml:"name"`
			Weight int    `yaml:"weight"`
		} `yaml:"to"`
		Port struct {
			TargetPort string `yaml:"targetPort"`
		} `yaml:"port"`
		TLS struct {
			Termination                   string `yaml:"termination"`
			InsecureEdgeTerminationPolicy string `yaml:"insecureEdgeTerminationPolicy"`
		} `yaml:"tls"`
		WildcardPolicy string `yaml:"wildcardPolicy"`
	} `yaml:"spec"`
}
