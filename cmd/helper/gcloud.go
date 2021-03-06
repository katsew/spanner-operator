package helper

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func GetGCPDefaults() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	gcloudConfigPath := fmt.Sprintf("%s/.config/gcloud", home)
	activeConfigPath := fmt.Sprintf("%s/active_config", gcloudConfigPath)
	b, err := ioutil.ReadFile(activeConfigPath)
	if err != nil {
		return "", "", err
	}
	configlFilePath := fmt.Sprintf("%s/configurations/config_%s", gcloudConfigPath, b)
	f, err := os.OpenFile(configlFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return "", "", err
	}
	var projectId string
	var instanceConfig string
	rd := bufio.NewReader(f)
	for {
		l, _, err := rd.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", err
		}
		s := string(l)
		if strings.HasPrefix(s, "project = ") {
			projectId = strings.TrimPrefix(s, "project = ")
		}
		if strings.HasPrefix(s, "region = ") {
			instanceConfig = fmt.Sprintf("regional-%s", strings.Trim(s, "region = "))
			if FindInstanceConfigByName(instanceConfig) == InstanceConfigUndefined {
				instanceConfig = GetRandomInstanceConfig()
			}
		}
	}
	return projectId, instanceConfig, nil
}
