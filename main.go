package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	cmdroot "github.com/grokspawn/stencil/cmd"
	registrylib "github.com/operator-framework/operator-registry/pkg/registry"
)

func main() {
	cmd := cmdroot.NewCmd()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := cmd.ExecuteContext(ctx); err != nil {
		agg, ok := err.(utilerrors.Aggregate)
		if !ok {
			os.Exit(1)
		}
		for _, e := range agg.Errors() {
			if _, ok := e.(registrylib.BundleImageAlreadyAddedErr); ok {
				os.Exit(2)
			}
			if _, ok := e.(registrylib.PackageVersionAlreadyAddedErr); ok {
				os.Exit(3)
			}
		}
		os.Exit(1)
	}
}
