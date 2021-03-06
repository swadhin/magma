/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// This starts the home subscriber server (hss) service.
package main

import (
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.MOCK_HSS)
	if err != nil {
		log.Fatalf("Error creating hss service: %s", err)
	}
	config, err := servicers.GetHSSConfig()
	if err != nil {
		log.Fatalf("Error getting hss config: %s", err)
	}
	servicer, err := servicers.NewHomeSubscriberServer(storage.NewMemorySubscriberStore(), config)
	if err != nil {
		log.Fatalf("Error creating home subscriber server: %s", err)
	}
	protos.RegisterHSSConfiguratorServer(srv.GrpcServer, servicer)

	// Start diameter server
	go func() {
		glog.V(2).Info("Starting home subscriber server")
		err := servicer.Start() // blocks
		glog.Error(err)
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running hss service: %s", err)
	}
}
