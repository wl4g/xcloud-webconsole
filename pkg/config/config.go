/**
 * Copyright 2017 ~ 2025 the original author or authors<Wanglsir@gmail.com, 983708408@qq.com>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package config

import (
	"fmt"

	utils "xcloud-webconsole/pkg/utils"
)

// GlobalProperties ...
type GlobalProperties struct {
	Admin      AdminProperties      `yaml:"admin"`
	Server     ServerProperties     `yaml:"server"`
	DataSource DataSourceProperties `yaml:"datasource"`
	Logging    LoggingProperties    `yaml:"logging"`
}

var (
	// GlobalConfig ...
	GlobalConfig GlobalProperties
)

// InitGlobalConfig global config properties.
func InitGlobalConfig(path string) {
	// Check configuraion path
	if !utils.ExistsFileOrDir(path) {
		path = defaultConfPath // fallback use default
		if !utils.ExistsFileOrDir(path) {
			panic(fmt.Sprintf("No such configuration file '%s'", path))
		}
	}
	fmt.Printf("Load configuration file '%s'", path)

	// Create default config.
	globalConfig := createDefaultProperties()

	c := utils.NewViperConfigurer()
	c.SetConfig(defaultWebConsoleYamlConfigContent, path)
	c.SetEnvPrefix("WEBCONSOLE") // Sets auto use env variables prefix

	// Load & parse
	if err := c.Parse(globalConfig); err != nil {
		panic(err)
	}

	// conf, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Read config '%s' error! %s", path, err)
	// 	panic(err)
	// }

	// err = yaml.Unmarshal(conf, globalConfig)
	// if err != nil {
	// 	fmt.Printf("Unmarshal config '%s' error! %s", path, err)
	// 	panic(err)
	// }

	// Post properties.
	afterPropertiesSet(globalConfig)

	// Sets Global configuration
	GlobalConfig = *globalConfig
}

// Create default configuration properties.
func createDefaultProperties() *GlobalProperties {
	return &GlobalProperties{}
}

// MetricExclude settings after initialization
func afterPropertiesSet(globalConfig *GlobalProperties) {
}

// RefreshConfig Refresh global config.
func RefreshConfig(config *GlobalProperties) {
	utils.CopyProperties(&config, &GlobalConfig)
}

const (
	defaultConfPath = "/etc/webconsole.yml"
)
