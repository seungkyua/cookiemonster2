package main

import (
	"fmt"
	"time"

	"github.com/seungkyua/cookiemonster2/src/domain"
)

func main() {
	config := domain.GetConfig()
	err := config.ReadConfig("../../config/")
	if err != nil {
		fmt.Println(err)
	}

	pm := &domain.PodManage{
		Started: false,
	}

	ticker := time.NewTicker(time.Second * time.Duration(config.Namespace[0].Resource[0].Interval))

	err = pm.MainLoop(config)
	if err != nil {
		panic(err.Error())
	}

	for range ticker.C {
		err = pm.MainLoop(config)
		if err != nil {
			panic(err.Error())
		}
	}
}
