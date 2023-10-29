package k8s_client

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/PoteeDev/k8s-controller/internal/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func ListPods(namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	fmt.Println("Get Kubernetes Pods")
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v", err)
		return nil, err
	}
	return pods, nil
}

func GetNamespaceInfo(namespace string) (*models.StandInfo, error) {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	pods, err := ListPods(namespace, clientset)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var stand models.StandInfo
	for _, pod := range pods.Items {
		var containers []models.Container
		for _, container := range pod.Status.ContainerStatuses {
			containers = append(containers, models.Container{
				Name:  container.Name,
				Image: container.Image,
				Ready: container.Ready,
			})
		}
		stand.Components = append(stand.Components, models.ComponentInfo{
			Name:       pod.Name,
			Containers: containers,
			Address:    pod.Status.PodIP,
			Status:     string(pod.Status.Phase),
		})
	}

	stand.TotalComponents = len(pods.Items)
	return &stand, err
}
