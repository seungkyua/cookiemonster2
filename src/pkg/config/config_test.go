package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/seungkyua/cookiemonster2/src/pkg/config"
	"gopkg.in/yaml.v2"
)

var (
	rawConfig = []byte(`
namespace:
  - name: openstack
    resource:
    - kind: deployment
      name: rabbitmq
      target: 1
      interval: 60
      duration: 600
      slack: true
    - kind: deployment
      name:
      target: 1
      interval: 60
      duration: 600
      slack: true
    - kind: daemonset
      name:
      target: 1
      interval: 60
      duration: 600
      slack: true
    - kind: statefulset
      name: mariadb
      target: 1
      interval: 60
      duration: 600
      slack: true

`)
)

func TestReadConfig(t *testing.T) {
	path := "../../../config/"
	var wantConfig config.Config
	if err := yaml.Unmarshal(rawConfig, &wantConfig); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(wantConfig)

	gotConfig, err := config.ReadConfig(path)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(gotConfig)

	if !reflect.DeepEqual(wantConfig, gotConfig) {
		t.Errorf("It was incorrect, got: %v, want: %v.", gotConfig, wantConfig)
	}

}
