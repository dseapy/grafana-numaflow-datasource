package query

// This file is based on "handler.go" in numaflow server, duplicated
// because fields aren't exported to be used outside of server package

import (
	"context"
	"errors"
	"fmt"
	"os"

	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"github.com/numaproj/numaflow/pkg/apis/proto/daemon"
	dfv1versiond "github.com/numaproj/numaflow/pkg/client/clientset/versioned"
	dfv1clients "github.com/numaproj/numaflow/pkg/client/clientset/versioned/typed/numaflow/v1alpha1"
	daemonclient "github.com/numaproj/numaflow/pkg/daemon/client"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsversiond "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	kubeClient     kubernetes.Interface
	metricsClient  *metricsversiond.Clientset
	numaflowClient dfv1clients.NumaflowV1alpha1Interface
	listOptions    metav1.ListOptions
}

func NewClient() (*Client, error) {
	var restConfig *rest.Config
	var err error
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = home + "/.kube/config"
		if _, err := os.Stat(kubeconfig); err != nil && os.IsNotExist(err) {
			kubeconfig = ""
		}
	}
	if kubeconfig != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		restConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig, %w", err)
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeclient, %w", err)
	}
	metricsClient := metricsversiond.NewForConfigOrDie(restConfig)
	numaflowClient := dfv1versiond.NewForConfigOrDie(restConfig).NumaflowV1alpha1()
	return &Client{
		kubeClient:     kubeClient,
		metricsClient:  metricsClient,
		numaflowClient: numaflowClient,
		// for now hard-code default limit, in future can allow overriding in data source or in each data query
		listOptions: metav1.ListOptions{Limit: 1000},
	}, nil
}

func (c *Client) ListNamespacesWithPipelines(ns string) ([]string, error) {
	l, err := c.numaflowClient.Pipelines(ns).List(context.Background(), c.listOptions)
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool)
	for _, pl := range l.Items {
		m[pl.Namespace] = true
	}
	namespaces := []string{}
	for k := range m {
		namespaces = append(namespaces, k)
	}
	return namespaces, nil
}

func (c *Client) ListNamespacesWithVertices(ns string) ([]string, error) {
	l, err := c.numaflowClient.Vertices(ns).List(context.Background(), c.listOptions)
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool)
	for _, pl := range l.Items {
		m[pl.Namespace] = true
	}
	namespaces := []string{}
	for k := range m {
		namespaces = append(namespaces, k)
	}
	return namespaces, nil
}

func (c *Client) ListNamespacesWithInterStepBufferServices(ns string) ([]string, error) {
	l, err := c.numaflowClient.InterStepBufferServices(ns).List(context.Background(), c.listOptions)
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool)
	for _, pl := range l.Items {
		m[pl.Namespace] = true
	}
	namespaces := []string{}
	for k := range m {
		namespaces = append(namespaces, k)
	}
	return namespaces, nil
}

func (c *Client) ListPipelines(ns string) ([]dfv1.Pipeline, error) {
	plList, err := c.numaflowClient.Pipelines(ns).List(context.Background(), c.listOptions)
	if err != nil {
		return nil, err
	}
	return plList.Items, nil
}

func (c *Client) ListVertices(ns string) ([]dfv1.Vertex, error) {
	vertexList, err := c.numaflowClient.Vertices(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return vertexList.Items, nil
}

func (c *Client) ListPipelineVertices(ns, pipeline string) ([]dfv1.Vertex, error) {
	lo := c.listOptions.DeepCopy()
	lo.LabelSelector = fmt.Sprintf("%s=%s", dfv1.KeyPipelineName, pipeline)
	vertices, err := c.numaflowClient.Vertices(ns).List(context.Background(), *lo)
	if err != nil {
		return nil, err
	}
	return vertices.Items, nil
}

func (c *Client) ListInterStepBufferServices(ns string) ([]dfv1.InterStepBufferService, error) {
	isbsvcList, err := c.numaflowClient.InterStepBufferServices(ns).List(context.Background(), c.listOptions)
	if err != nil {
		return nil, err
	}
	return isbsvcList.Items, nil
}

func (c *Client) GetPipeline(ns, pipeline string) (*dfv1.Pipeline, error) {
	pl, err := c.numaflowClient.Pipelines(ns).Get(context.Background(), pipeline, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func (c *Client) GetPipelineVertex(ns, pipeline, vertex string) (*dfv1.Vertex, error) {
	lo := c.listOptions.DeepCopy()
	lo.LabelSelector = fmt.Sprintf("%s=%s,%s=%s", dfv1.KeyPipelineName, pipeline, dfv1.KeyVertexName, vertex)
	vertices, err := c.numaflowClient.Vertices(ns).List(context.Background(), *lo)
	if err != nil {
		return nil, err
	}
	if len(vertices.Items) == 0 {
		return nil, errors.New(fmt.Sprintf("Vertex %q not found", vertex))
	}
	return &vertices.Items[0], err
}

func (c *Client) GetInterStepBufferService(ns, isbsvc string) (*dfv1.InterStepBufferService, error) {
	i, err := c.numaflowClient.InterStepBufferServices(ns).Get(context.Background(), isbsvc, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (c *Client) ListVertexPods(ns, pipeline, vertex string) ([]v1.Pod, error) {
	lo := c.listOptions.DeepCopy()
	lo.LabelSelector = fmt.Sprintf("%s=%s,%s=%s", dfv1.KeyPipelineName, pipeline, dfv1.KeyVertexName, vertex)
	pods, err := c.kubeClient.CoreV1().Pods(ns).List(context.Background(), *lo)
	if err != nil {
		return nil, err
	}
	return pods.Items, err
}

func (c *Client) ListPodsMetrics(ns string) ([]v1beta1.PodMetrics, error) {
	l, err := c.metricsClient.MetricsV1beta1().PodMetricses(ns).List(context.Background(), c.listOptions)
	if err != nil {
		return nil, err
	}
	return l.Items, nil
}

func (c *Client) GetPodMetrics(ns, po string) (*v1beta1.PodMetrics, error) {
	m, err := c.metricsClient.MetricsV1beta1().PodMetricses(ns).Get(context.Background(), po, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Client) ListPipelineEdges(ns, pipeline string) ([]*daemon.BufferInfo, error) {
	client, err := daemonclient.NewDaemonServiceClient(daemonSvcAddress(ns, pipeline))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()
	l, err := client.ListPipelineBuffers(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (c *Client) GetPipelineEdge(ns, pipeline, edge string) (*daemon.BufferInfo, error) {
	client, err := daemonclient.NewDaemonServiceClient(daemonSvcAddress(ns, pipeline))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()
	// Assume edge is the buffer name
	i, err := client.GetPipelineBuffer(context.Background(), pipeline, edge)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (c *Client) GetVertexMetrics(ns, pipeline, vertex string) (*daemon.VertexMetrics, error) {
	client, err := daemonclient.NewDaemonServiceClient(daemonSvcAddress(ns, pipeline))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()
	l, err := client.GetVertexMetrics(context.Background(), pipeline, vertex)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (c *Client) GetVertexWatermark(ns, pipeline, vertex string) (*daemon.VertexWatermark, error) {
	client, err := daemonclient.NewDaemonServiceClient(daemonSvcAddress(ns, pipeline))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()
	l, err := client.GetVertexWatermark(context.Background(), pipeline, vertex)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func daemonSvcAddress(ns, pipeline string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", fmt.Sprintf("%s-daemon-svc", pipeline), ns, dfv1.DaemonServicePort)
}
