package pkg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/linuxsuren/github-proxy/pkg/module"
)

func ParsePipelineTask(taskPath string) (task *module.PipelineTask, err error) {
	var buf []byte

	buf, err = ioutil.ReadFile(taskPath)
	if err != nil {
		return
	}

	task = &module.PipelineTask{}
	err = yaml.Unmarshal(buf, task)

	return
}