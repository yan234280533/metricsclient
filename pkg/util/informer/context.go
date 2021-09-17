package util

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

type Client struct {
	Master                  string
	Kubeconfig              string
	NodeName                string
	kubeClient              clientset.Interface
	nodeFactory, podFactory informers.SharedInformerFactory
}

const (
	nodeNameField      = "metadata.name"
	specNodeNameField  = "spec.nodeName"
	statusPhaseFiled   = "status.phase"
	informerSyncPeriod = time.Minute
)

func (c *Client) lazyInit() {
	if c.kubeClient != nil {
		return
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags(c.Master, c.Kubeconfig)
	if err != nil {
		klog.Warning(err)
		klog.Warning("fall back to creating fake kube-client")
		// create a fake client to test caelus without k8s
		c.kubeClient = fake.NewSimpleClientset()
	} else {
		c.kubeClient = clientset.NewForConfigOrDie(kubeconfig)
	}
}

// GetKubeClient returns k8s client interface
func (c *Client) GetKubeClient() clientset.Interface {
	c.lazyInit()
	return c.kubeClient
}

// GetPodFactory returns pod factory
func (c *Client) GetPodFactory() informers.SharedInformerFactory {
	if c.podFactory == nil {
		c.podFactory = informers.NewSharedInformerFactoryWithOptions(c.GetKubeClient(), informerSyncPeriod,
			informers.WithTweakListOptions(func(options *metav1.ListOptions) {
				options.FieldSelector = fields.AndSelectors(fields.OneTermEqualSelector(specNodeNameField, c.NodeName),
					fields.OneTermNotEqualSelector(statusPhaseFiled, "Succeeded"),
					fields.OneTermNotEqualSelector(statusPhaseFiled, "Failed")).String()
			}))
	}
	return c.podFactory
}

// GetNodeFactory returns node factory
func (c *Client) GetNodeFactory() informers.SharedInformerFactory {
	if c.nodeFactory == nil {
		c.nodeFactory = informers.NewSharedInformerFactoryWithOptions(c.GetKubeClient(), informerSyncPeriod,
			informers.WithTweakListOptions(func(options *metav1.ListOptions) {
				options.FieldSelector = fields.OneTermEqualSelector(nodeNameField, c.NodeName).String()
			}))
	}
	return c.nodeFactory
}

// Run starts k8s informers
func (c *Client) Run(stop <-chan struct{}) {
	if c.podFactory != nil {
		c.podFactory.Start(stop)
		c.podFactory.WaitForCacheSync(stop)
	}

	if c.nodeFactory != nil {
		c.nodeFactory.Start(stop)
		c.nodeFactory.WaitForCacheSync(stop)
	}
}
