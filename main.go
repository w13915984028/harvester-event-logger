//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go
//go:generate /bin/bash scripts/generate-manifest

package main

import (
  "encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

  "k8s.io/client-go/tools/cache"

	"github.com/w13915984028/harvester-event-logger/pkg/controller/eventlogger"
	"github.com/w13915984028/harvester-event-logger/pkg/config"
)

var (
	VERSION = "v0.0.1"
)


func main() {
	app := cli.NewApp()
	app.Name = "harvester-event-logger"
	app.Usage = "harvester-event-logger logs all events in the harvester cluster"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kubeconfig, k",
			EnvVar: "KUBECONFIG",
			Value:  "",
			Usage:  "Kubernetes config files, e.g. $HOME/.kube/config",
		},
		cli.StringFlag{
			Name:   "master, m",
			EnvVar: "MASTERURL",
			Value:  "",
			Usage:  "Kubernetes cluster master URL.",
		},
	}
	app.Action = func(c *cli.Context) {
		if err := run(c); err != nil {
			panic(err)
		}
	}

	if err := app.Run(os.Args); err != nil {
		klog.Error(err)
	}
}

func run(c *cli.Context) error {
	masterURL := c.String("master")
	kubeconfig := c.String("kubeconfig")

  ctx := signals.SetupSignalContext()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		return fmt.Errorf("error building config from flags: %w", err)
	}

	options := &config.Options{

	}

	management, err := config.SetupManagement(ctx, cfg, options)
	if err != nil {
		klog.Fatalf("Error building harvester controllers: %s", err.Error())
	}

	callback := func(ctx context.Context) {
		if err := management.Register(ctx, cfg, config.RegisterFuncList); err != nil {
			panic(err)
		}

		if err := management.Start(threadiness); err != nil {
			panic(err)
		}

		<-ctx.Done()
	}

	if leaderelection {
		leader.RunOrDie(ctx, "harvester-system", "harvester-event-logger", client, callback)
	} else {
		callback(ctx)
	}

	return nil
}
