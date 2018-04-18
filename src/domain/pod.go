package domain

import (
	"log"
	"math/rand"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apps_v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
	"k8s.io/api/core/v1"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/apis/core/pods"
)

var ron = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomInt(i int) int {
	return ron.Intn(i)
}

type PodManage struct {
	Ctx		context.Context
	Cancel	context.CancelFunc
	Started	bool
}

func connect() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return clientset, nil
}

func (m *PodManage) Start(c *Config) error {
	log.Printf("Cookie Time!!! Random feast starting")

	tick := time.Tick(time.Duration(c.Namespace[0].Resource[0].Interval) * time.Second)
	for {
		select {
		case <-m.Ctx.Done():
			return nil
		case <-tick:
			err := m.mainLoop(c)
			if err != nil {
				return err
			}
		}
	}

	//for _, ns := range c.Namespace {
	//	for _, res := range ns.Resource {
	//		log.Printf("Cookie Time!!! Random feast starting on %s in namespace %s", res.Kind, ns)
	//		tick := time.Tick(time.Duration(res.Interval) * time.Second)
	//		err := m.mainLoop(tick, res.Kind)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}
	return nil
}

func (m *PodManage) mainLoop(c *Config) error {
	for _, ns := range c.Namespace {
		for _, res := range ns.Resource {
			pod, found, err := m.SelectVictimPod(c, ns.Name, res.Kind)
			if err != nil {
				return err
			} else if found {
				go killPod(pod, ns.Name, res.Kind)
			}
		}
	}
	return nil
}

func (m *PodManage) Stop(c *Config) {
	defer m.Cancel()
	m.Started = false
}

func (m *PodManage) SelectVictimPod(c *Config, ns string, kind string) (*v1.Pod, bool, error) {
	var found bool
	var pod v1.Pod

	con, err := connect()
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	//x := randomInt(len(c.Namespace))
	//pods, err := clientset.AppsV1().Deployments(ns).List(metaV1.ListOptions{})
	//pods, err := clientset.CoreV1().Pods(c.Namespace[x].Name).List(metaV1.ListOptions{})

	// query named deployment set via name provided
	lo := metaV1.ListOptions{FieldSelector: "kind=" + kind}

	pods, err := con.CoreV1().Pods(ns).List(lo)
	if err != nil {
		log.Println(err)
		return nil, false, err
	} else if len(pods.Items) > 0 {
		found = true
		item := randomInt(len(pods.Items))
		pod = pods.Items[item]
		log.Printf("Found %s %s in namespace %s\n", kind, pod.Name, ns)
	} else {
		log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, pod.Name, ns)
	}

	return &pod, found, nil
}

func killPod(pod *v1.Pod, ns string, kind string) error {
	con, err := connect()
	if err != nil {
		log.Println(err)
		return err
	}

	err = con.CoreV1().Pods(ns).Delete(pod.Name, &metaV1.DeleteOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Eating pod %s NOM NOM NOM!!!!", pod.Name)

	return nil
}
