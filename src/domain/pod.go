package domain

import (
	"log"
	"math/rand"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

var randum = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomInt(i int) int {
	return randum.Intn(i)
}

type VictimPod struct {
	Name string
	Kind string
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

func (v *VictimPod) GetVictimPod(config *Config) error {
	clientset, err := connect()
	if err != nil {
		log.Println(err)
		return err
	}

	x := randomInt(len(config.Namespace))

	pods, err := clientset.CoreV1().Pods(config.Namespace[x].Name).List(metaV1.ListOptions{})
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(*pods)
	return nil
}
