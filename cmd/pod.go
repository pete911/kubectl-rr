package cmd

import (
	"fmt"
	"github.com/pete911/kubectl-rr/internal"
	"github.com/pete911/kubectl-rr/internal/out"
	"github.com/spf13/cobra"
	"os"
)

var (
	cmdList = &cobra.Command{
		Use:   "pod",
		Short: "pod resource recommendation",
		Long:  "",
		RunE:  runPodCmd,
	}
	podFlags PodFlags
)

func init() {
	RootCmd.AddCommand(cmdList)
	InitPodFlags(cmdList, &podFlags)
}

func runPodCmd(_ *cobra.Command, args []string) error {
	cfg := podFlags.ToConfig(args)
	pods, err := internal.GetPods(RestConfig(), cfg)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	out.PrintPods(pods)
	return nil
}
