package scenario

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

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

// NFClients and this file is based on handler.go in numaflow server
type NFClients struct {
	kubeClient     kubernetes.Interface
	metricsClient  *metricsversiond.Clientset
	numaflowClient dfv1clients.NumaflowV1alpha1Interface
}

// NewNFClients creates various clients used to get data about numaflow from various endpoints
func NewNFClients() (*NFClients, error) {
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
	return &NFClients{
		kubeClient:     kubeClient,
		metricsClient:  metricsClient,
		numaflowClient: numaflowClient,
	}, nil
}

// ListPipelines is used to provide all the numaflow pipelines in a given namespace
func (c *NFClients) ListPipelines(ns string) ([]dfv1.Pipeline, error) {
	plList, err := c.numaflowClient.Pipelines(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return plList.Items, nil
}

// GetPipeline is used to provide the spec of a given numaflow pipeline
func (c *NFClients) GetPipeline(ns, pipeline string) (*dfv1.Pipeline, error) {
	pl, err := c.numaflowClient.Pipelines(ns).Get(context.Background(), pipeline, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pl, nil
}

// ListNamespacesWithPipelines is used to provide all the namespaces that have numaflow pipelines running
func (c *NFClients) ListNamespacesWithPipelines(ns string) ([]string, error) {
	l, err := c.numaflowClient.Pipelines(ns).List(context.Background(), metav1.ListOptions{})
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

// ListNamespacesWithVertices is used to provide all the namespaces that have numaflow pipelines running
func (c *NFClients) ListNamespacesWithVertices(ns string) ([]string, error) {
	l, err := c.numaflowClient.Vertices(ns).List(context.Background(), metav1.ListOptions{})
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

// ListNamespacesWithInterStepBufferServices is used to provide all the namespaces that have numaflow pipelines running
func (c *NFClients) ListNamespacesWithInterStepBufferServices(ns string) ([]string, error) {
	l, err := c.numaflowClient.InterStepBufferServices(ns).List(context.Background(), metav1.ListOptions{})
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

// ListInterStepBufferServices is used to provide all the interstepbuffer services in a namespace
func (c *NFClients) ListInterStepBufferServices(ns string) ([]dfv1.InterStepBufferService, error) {
	isbsvcList, err := c.numaflowClient.InterStepBufferServices(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return isbsvcList.Items, nil
}

// GetInterStepBufferService is used to provide the spec of the interstep buffer service
func (c *NFClients) GetInterStepBufferService(ns, isbsvc string) (*dfv1.InterStepBufferService, error) {
	i, err := c.numaflowClient.InterStepBufferServices(ns).Get(context.Background(), isbsvc, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return i, nil
}

// ListVertices is used to provide all the vertices of a pipeline
func (c *NFClients) ListVertices(ns string) ([]dfv1.Vertex, error) {
	vertexList, err := c.numaflowClient.Vertices(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return vertexList.Items, nil
}

// ListPipelineVertices is used to provide all the vertices of a pipeline
func (c *NFClients) ListPipelineVertices(ns, pipeline string) ([]dfv1.Vertex, error) {
	vertices, err := c.numaflowClient.Vertices(ns).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", dfv1.KeyPipelineName, pipeline),
	})
	if err != nil {
		return nil, err
	}
	return vertices.Items, nil
}

// GetVertex is used to provide the vertex spec
func (c *NFClients) GetVertex(ns, vertex, pipeline string) (*dfv1.Vertex, error) {
	vertices, err := c.numaflowClient.Vertices(ns).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=%s", dfv1.KeyPipelineName, pipeline, dfv1.KeyVertexName, vertex),
	})
	if err != nil {
		return nil, err
	}
	if len(vertices.Items) == 0 {
		return nil, errors.New(fmt.Sprintf("Vertex %q not found", vertex))
	}
	return &vertices.Items[0], err
}

// ListVertexPods is used to provide all the pods of a vertex
func (c *NFClients) ListVertexPods(ns, pipeline, vertex, limit, cont string) ([]v1.Pod, error) {
	lmt, _ := strconv.ParseInt(limit, 10, 64)
	pods, err := c.kubeClient.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=%s", dfv1.KeyPipelineName, pipeline, dfv1.KeyVertexName, vertex),
		Limit:         lmt,
		Continue:      cont,
	})
	if err != nil {
		return nil, err
	}
	return pods.Items, err
}

// ListPodsMetrics is used to provide a list of all metrics in all the pods
func (c *NFClients) ListPodsMetrics(ns, limit, cont string) ([]v1beta1.PodMetrics, error) {
	lmt, _ := strconv.ParseInt(limit, 10, 64)
	l, err := c.metricsClient.MetricsV1beta1().PodMetricses(ns).List(context.Background(), metav1.ListOptions{
		Limit:    lmt,
		Continue: cont,
	})
	if err != nil {
		return nil, err
	}
	return l.Items, nil
}

// GetPodMetrics is used to provide the metrics like CPU/Memory utilization for a pod
func (c *NFClients) GetPodMetrics(ns, po string) (*v1beta1.PodMetrics, error) {
	m, err := c.metricsClient.MetricsV1beta1().PodMetricses(ns).Get(context.Background(), po, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListPipelineEdges is used to provide information about all the pipeline edges
func (c *NFClients) ListPipelineEdges(ns, pipeline string) ([]*daemon.BufferInfo, error) {
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

// GetPipelineEdge is used to provide information about a single pipeline edge
func (c *NFClients) GetPipelineEdge(ns, pipeline, edge string) (*daemon.BufferInfo, error) {
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

// GetVertexMetrics is used to provide information about the vertex including processing rates.
func (c *NFClients) GetVertexMetrics(ns, pipeline, vertex string) (*daemon.VertexMetrics, error) {
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

// GetVertexWatermark is used to provide the head watermark for a given vertex
func (c *NFClients) GetVertexWatermark(ns, pipeline, vertex string) (*daemon.VertexWatermark, error) {
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
