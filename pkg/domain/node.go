package domain

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type Namelist struct {
	ServerName string
	NodeName []string
}

func Returnnamelist() *Namelist {
	var l *Namelist
	l = &Namelist{ServerName:"Gyutae"}
	err := NodeList(l)
	if err != nil {
		log.Println(err)
	}
	return l
}

func NodeList(ls *Namelist) error {
	con, err := Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	var kubecore = con.CoreV1()
	var lo = metav1.ListOptions{}
	nodelist, err := kubecore.Nodes().List(lo)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, node := range nodelist.Items {
		ls.NodeName = append(ls.NodeName, node.Name)
	}
	return nil

}
