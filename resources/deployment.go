package resources

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"safaricom/pipeline/files"
	"safaricom/pipeline/models"
	"strings"

	"gopkg.in/yaml.v2"
)

func writeYAMLToFile(body interface{}, resource_type string, name string) (file string, err error) {
	filePath := "files/" + resource_type + "_" + name + ".yaml"

	var data string
	if reflect.TypeOf(body).Kind() == reflect.String {
		data = fmt.Sprintf("%v", body)
	} else {
		updatedYaml, er := yaml.Marshal(body)
		data = strings.ReplaceAll(string(updatedYaml), "service-name", name)

		if er != nil {
			return
		}
	}

	return files.WriteFile(filePath, string(data))
}

func generateHPA(request models.Payload) (file string, err error) {
	yamlFile, err := os.ReadFile("files/hpa.yaml")
	if err != nil {
		return
	}

	var hpa models.HorizontalPodAutoscaler
	err = yaml.Unmarshal(yamlFile, &hpa)
	if err != nil {
		return
	}

	hpa.Metadata.Name = request.ServiceName
	hpa.Metadata.Namespace = request.Namespace
	hpa.Spec.ScaleTargetRef.Name = request.ServiceName
	hpa.Spec.MinReplicas = request.MinPods
	hpa.Spec.MaxReplicas = request.MaxPods

	return writeYAMLToFile(hpa, "hpa", request.ServiceName)
}

func generateRoute(request models.Payload) (file string, err error) {
	yamlFile, err := os.ReadFile("files/route.yaml")
	if err != nil {
		return
	}

	var route models.Route
	err = yaml.Unmarshal(yamlFile, &route)
	if err != nil {
		return
	}

	route.Metadata.Name = request.ServiceName
	route.Metadata.Namespace = request.Namespace
	route.Spec.To.Name = request.ServiceName
	route.Spec.Host = fmt.Sprintf("%s.%s.apps.devocp.safaricom.net", request.ServiceName, request.Namespace)

	return writeYAMLToFile(route, "route", request.ServiceName)
}

func generateServiceFile(request models.Payload) (file string, err error) {
	yamlFile, err := os.ReadFile("files/service.yaml")
	if err != nil {
		return
	}

	var service models.Service
	err = yaml.Unmarshal(yamlFile, &service)
	if err != nil {
		return
	}

	service.Metadata.Name = request.ServiceName
	service.Metadata.Namespace = request.Namespace
	service.Metadata.Labels.App = request.ServiceName
	service.Spec.Selector.App = request.ServiceName
	service.Spec.Ports[0].Port = request.Port
	service.Spec.Ports[0].TargetPort = request.Port
	service.Spec.Ports[0].NodePort = request.NodePort

	return writeYAMLToFile(service, "service", request.ServiceName)
}

func generateImagePolicy(request models.Payload) (file string, err error) {
	yamlFile, err := os.ReadFile("files/imagePolicy.yaml")
	if err != nil {
		return
	}

	var policy models.ImagePolicy
	err = yaml.Unmarshal(yamlFile, &policy)
	if err != nil {
		return
	}

	policy.Metadata.Name = fmt.Sprintf("%s-policy", request.ServiceName)
	policy.Spec.ImageRepositoryRef.Name = fmt.Sprintf("%s-repo", request.ServiceName)
	policy.Metadata.Namespace = request.Namespace

	policy.Spec.FilterTags.Pattern = fmt.Sprintf("^%s", request.FluxTag)

	return writeYAMLToFile(policy, "policy", request.ServiceName)
}

func generateImageRepository(request models.Payload) (file string, err error) {
	yamlFile, err := os.ReadFile("files/imageRepo.yaml")
	if err != nil {
		return
	}

	var repo models.ImageRepository
	err = yaml.Unmarshal(yamlFile, &repo)
	if err != nil {
		return
	}

	end := strings.Index(request.Image, ":")

	if end > -1 {
		repo.Spec.Image = request.Image[0:end]

	} else {
		repo.Spec.Image = request.Image
		fmt.Println(request.ServiceName, ": ", request.Image)
	}

	repo.Metadata.Name = fmt.Sprintf("%s-repo", request.ServiceName)

	repo.Metadata.Namespace = request.Namespace
	repo.Spec.SecretRef.Name = request.ImagePullSecrets

	return writeYAMLToFile(repo, "repo", request.ServiceName)
}

func generateDeploymentFile(request models.Payload) (file string, err error) {
	resources := make([]string, 0)
	repo, err := generateImageRepository(request)
	if err != nil {
		return
	}
	defer files.DeleteFile(repo)
	repoYAML, err := os.ReadFile(repo)
	if err != nil {
		return
	}

	policy, err := generateImagePolicy(request)
	if err != nil {
		return
	}

	defer files.DeleteFile(policy)
	policyYAML, err := os.ReadFile(policy)
	if err != nil {
		return
	}

	resources = append(resources, string(repoYAML))
	resources = append(resources, string(policyYAML))

	yamlFile, err := os.ReadFile("files/deployment.yaml")
	if err != nil {
		return
	}

	var deployment models.Deployment
	err = yaml.Unmarshal(yamlFile, &deployment)
	if err != nil {
		return
	}

	deployment.Metadata.Name = request.ServiceName
	deployment.Metadata.Namespace = request.Namespace
	deployment.Spec.Selector.MatchLabels.App = request.ServiceName
	deployment.Spec.Template.Metadata.Labels.App = request.ServiceName

	if request.MaxPods <= 1 {
		deployment.Spec.Replicas = request.MinPods
	}

	deployment.Spec.Template.Spec.Containers[0].Name = request.ServiceName
	deployment.Spec.Template.Spec.Containers[0].Image = request.Image
	deployment.Spec.Template.Spec.ImagePullSecrets[0].Name = request.ImagePullSecrets

	if request.Port > 0 {
		deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = request.Port
	}

	if request.Resources != (models.Resources{}) {
		err = updateResources(&deployment, request)
		if err != nil {
			return
		}
	}

	if request.Env != nil {
		err = updateEnviromentVariables(&deployment, request)
		if err != nil {
			return
		}
	}

	deploymentFile, err := writeYAMLToFile(deployment, "deployment-resource", request.ServiceName)

	defer files.DeleteFile(deploymentFile)
	deploymentYAML, err := os.ReadFile(deploymentFile)
	if err != nil {
		return
	}

	image := fmt.Sprintf("%s # {\"$imagepolicy\":\"%s:%s-policy\"}", request.Image, request.Namespace, request.ServiceName)
	resources = append(resources, strings.Replace(string(deploymentYAML), request.Image, image, 1))

	object := strings.Join(resources, "---\n")

	return writeYAMLToFile(object, "deployment", request.ServiceName)
}

func getAll(request models.Payload) (file string, err error) {
	resources := make([]string, 0)

	deployment, err := generateDeploymentFile(request)
	if err != nil {
		return
	}

	defer files.DeleteFile(deployment)
	deploymentYAML, err := os.ReadFile(deployment)
	if err != nil {
		return
	}

	resources = append(resources, string(deploymentYAML))

	if request.Port > 0 {
		service, errr := generateServiceFile(request)
		defer files.DeleteFile(service)
		if errr != nil {
			return "", errr
		}

		serviceYAML, errr := os.ReadFile(service)
		if errr != nil {
			return "", errr
		}
		resources = append(resources, string(serviceYAML))

		//Generate route
		route, errr := generateRoute(request)
		defer files.DeleteFile(route)
		if errr != nil {
			return "", errr
		}

		routeYAML, errr := os.ReadFile(route)
		if errr != nil {
			return "", errr
		}
		resources = append(resources, string(routeYAML))

	}

	if request.MaxPods > 1 {
		hpa, errr := generateHPA(request)
		if errr != nil {
			return "", err
		}

		defer files.DeleteFile(hpa)
		hpaYAML, errr := os.ReadFile(hpa)
		if errr != nil {
			return "", errr
		}
		resources = append(resources, string(hpaYAML))
	}

	object := strings.Join(resources, "---\n")
	return writeYAMLToFile(object, "all", request.ServiceName)
}

func GetResource(request models.Payload, resource string) (string, error) {

	if resource == "all" {
		return getAll(request)
	}

	if resource == "deployment" {
		return generateDeploymentFile(request)

	} else if resource == "service" {
		return generateServiceFile(request)
	} else if resource == "hpa" {
		return generateHPA(request)
	} else if resource == "route" {
		return generateRoute(request)
	}

	return "", errors.New("Resource type {" + resource + "} not supported")
}

func updateEnviromentVariables(deployment *models.Deployment, request models.Payload) error {
	deployment.Spec.Template.Spec.Containers[0].Env = request.Env
	return nil
}

func updateResources(deployment *models.Deployment, request models.Payload) (err error) {
	deployment.Spec.Template.Spec.Containers[0].Resources = request.Resources
	return nil
}
