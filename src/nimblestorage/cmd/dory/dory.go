/*
(c) Copyright 2017 Hewlett Packard Enterprise Development LP

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"nimblestorage/pkg/jconfig"
	flexvol "nimblestorage/pkg/k8s/flexvol"
	"nimblestorage/pkg/util"
)

var (
	// Version contains the current version added by the build process
	Version = "dev"
	// Commit contains the commit id added by the build process
	Commit = "unknown"

	dockerVolumePluginSocketPath = "/run/docker/plugins/nimble.sock"
	stripK8sFromOptions          = true
	logFilePath                  = "/var/log/dory.log"
	debug                        = false
	createVolumes                = true
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough args")
		return
	}

	overridden := initialize()
	util.OpenLogFile(logFilePath, 10, 4, 90, debug)
	defer util.CloseLogFile()
	pid := os.Getpid()
	util.LogInfo.Printf("[%d] entry  : Driver=%s Version=%s-%s Socket=%s Overridden=%t", pid, filepath.Base(os.Args[0]), Version, Commit, dockerVolumePluginSocketPath, overridden)

	driverCommand := os.Args[1]
	util.LogInfo.Printf("[%d] request: %s %v", pid, driverCommand, os.Args[2:])
	flexvol.Config(dockerVolumePluginSocketPath, stripK8sFromOptions, createVolumes)
	mess := flexvol.Handle(driverCommand, os.Args[2:])
	util.LogInfo.Printf("[%d] reply  : %s %v: %v", pid, driverCommand, os.Args[2:], mess)

	fmt.Println(mess)
}

func initialize() bool {
	override := false

	// don't log anything in initialize because we haven't open a log file yet.
	c, err := jconfig.NewConfig(fmt.Sprintf("%s%s", os.Args[0], ".json"))
	if err != nil {
		return false
	}
	s, err := c.GetStringWithError("logFilePath")
	if err == nil && s != "" {
		override = true
		logFilePath = s
	}
	s, err = c.GetStringWithError("dockerVolumePluginSocketPath")
	if err == nil && s != "" {
		override = true
		dockerVolumePluginSocketPath = s
	}
	b, err := c.GetBool("logDebug")
	if err == nil {
		override = true
		debug = b
	}
	b, err = c.GetBool("stripK8sFromOptions")
	if err == nil {
		override = true
		stripK8sFromOptions = b
	}
	b, err = c.GetBool("createVolumes")
	if err == nil {
		override = true
		createVolumes = b
	}
	return override
}
