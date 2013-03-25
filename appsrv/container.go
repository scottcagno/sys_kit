// ------------
// container.go ::: application container
// ------------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package appsrv

import (
	"fmt"
	"os"
	"time"
)

var (
	null, _ = os.Create(os.DevNull)
	logf, _ = os.Create("main.log")
)

type Container struct {
	ContRoot string
	exitChan chan string
	ProcAttr os.ProcAttr
}

func NewContainer(contRoot string) *Container {
	self := &Container{
		ContRoot: contRoot,
		exitChan: make(chan string, 0),
	}
	self.ProcAttr.Files = []*os.File{null, null, logf}
	return self
}

func (self *Container) Run(appName string) *Container {
	go self.Spawn(appName)
	go func() {
		for {
			select {
			case app := <-self.exitChan:
				LogString("[%v]\t%q EXITED, RESPAWNING\n", time.Now(), app)
				go self.Spawn(app)
			}
		}
	}()
	return self
}

func (self *Container) Spawn(appName string) {
	proc, err := os.StartProcess(self.ContRoot+"/"+appName+"/"+appName, nil, &self.ProcAttr)
	if err != nil {
		LogError(err)
	}
	procState, _ := proc.Wait()
	if procState.Exited() {
		self.exitChan <- appName
	}
}

func LogString(fs string, args ...interface{}) {
	fmt.Fprintf(logf, fs, args...)
}

func LogError(err error) {
	fmt.Fprintf(logf, "[%v]\t%v\n", time.Now(), err)
}
