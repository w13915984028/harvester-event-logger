package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rancher/wrangler/pkg/leader"
	"github.com/rancher/wrangler/pkg/signals"

	"github.com/urfave/cli"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//	"k8s.io/klog"

	"github.com/w13915984028/harvester-event-logger/pkg/config"
	"github.com/w13915984028/harvester-event-logger/pkg/controller"
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
		panic(err)
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

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("Error get client from kubeconfig: %w", err)
	}

	options := &config.Options{}

	management, err := config.SetupManagement(ctx, cfg, options)
	if err != nil {
		return fmt.Errorf("Error building harvester controllers: %w", err)
	}

	callback := func(ctx context.Context) {
		if err := management.Register(ctx, cfg, controller.RegisterFuncList); err != nil {
			panic(err)
		}

		if err := management.Start(2); err != nil {
			panic(err)
		}

		<-ctx.Done()
	}

	// TBD, when multi instances are deployed, the leaderelection is required
	//  currently, only one instance is deployed
	//if leaderelection {
	if false {
		leader.RunOrDie(ctx, "cattle-logging-system", "harvester-event-logger", client, callback)
	} else {
		callback(ctx)
	}

	return nil
}
