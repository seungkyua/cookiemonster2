package domain

import (
	"fmt"
	"log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gofish "github.com/stmcginnis/gofish"
	redfish "github.com/stmcginnis/gofish/redfish"
	"time"
)

type Namelist struct {
	ServerName string
	NodeName []string
}

func Returnnamelist() *Namelist {
	var l *Namelist
	l = &Namelist{}
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

func Reboot() string {
	path := "../config"
	config = &Config{}
	if err := config.ReadConfig(path); err != nil {
		log.Println(err)
		return "Fail to restart node"
	}
	// Create a new instance of gofish client, ignoring self-signed certs
	fish_config := gofish.ClientConfig{
		Endpoint: config.Bmcad,
		Username: "my-username",
		Password: "my-password",
		Insecure: true,
	}

	c, err := gofish.Connect(fish_config)
	if err != nil {
		panic(err)
		return "Fail to restart node"
	}
	defer c.Logout()

	// Attached the client to service root
	service := c.Service

	// Query the computer systems
	ss, err := service.Systems()
	if err != nil {
		panic(err)
		return "Fail to restart node"
	}

	// Creates a boot override to pxe once
	bootOverride := redfish.Boot{
		BootSourceOverrideTarget:  redfish.PxeBootSourceOverrideTarget,
		BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	}
	go func(){
		for _, system := range ss {
			fmt.Printf("System: %#v\n\n", system)
			err := system.SetBoot(bootOverride)
			if err != nil {
				panic(err)
			}
			err = system.Reset(redfish.ForceRestartResetType)
			if err != nil {
				panic(err)
			}
			time.Sleep(1000)
		}
	}()
	return "Successfully restart node"
}
