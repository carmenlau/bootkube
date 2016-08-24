package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/util"

	"github.com/coreos/bootkube/pkg/bootkube"
)

var (
	cmdStart = &cobra.Command{
		Use:          "start",
		Short:        "Start the bootkube service",
		Long:         "",
		PreRunE:      validateStartOpts,
		RunE:         runCmdStart,
		SilenceUsage: true,
	}

	startOpts struct {
		assetDir       string
		etcdServer     string
		clusterIPRange string
	}
)

func init() {
	cmdRoot.AddCommand(cmdStart)
	cmdStart.Flags().StringVar(&startOpts.etcdServer, "etcd-server", "http://127.0.0.1:2379", "Single etcd node to use during bootkube bootstrap process.")
	cmdStart.Flags().StringVar(&startOpts.assetDir, "asset-dir", "", "Path to the cluster asset directory. Expected layout genereted by the `bootkube render` command.")
	cmdStart.Flags().StringVar(&startOpts.clusterIPRange, "cluster-ip-range", "10.3.0.0/24", "The CIDR of Service Cluster IP Range")
}

func runCmdStart(cmd *cobra.Command, args []string) error {
	etcdServer, err := url.Parse(startOpts.etcdServer)
	if err != nil {
		return fmt.Errorf("Invalid etcd etcdServer %q: %v", startOpts.etcdServer, err)
	}

	bk, err := bootkube.NewBootkube(bootkube.Config{
		AssetDir:       startOpts.assetDir,
		EtcdServer:     etcdServer,
		ClusterIPRange: startOpts.clusterIPRange,
	})

	if err != nil {
		return err
	}

	util.InitLogs()
	defer util.FlushLogs()

	return bk.Run()
}

func validateStartOpts(cmd *cobra.Command, args []string) error {
	if startOpts.etcdServer == "" {
		return errors.New("missing required flag: --etcd-server")
	}
	if startOpts.assetDir == "" {
		return errors.New("missing required flag: --asset-dir")
	}
	return nil
}
