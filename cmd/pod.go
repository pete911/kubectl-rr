package cmd

import (
	"fmt"
	"github.com/pete911/kubectl-rr/internal"
	"github.com/pete911/kubectl-rr/internal/out"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
	// no namespace means all namespaces
	if podFlags.AllNamespaces {
		podFlags.Namespace = ""
	}

	// additional arguments are considered to be pod names, add to field selector flags
	for _, v := range args {
		fieldSelectors := strings.Split(podFlags.FieldSelector, ",")
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.name=%s", v))
		podFlags.FieldSelector = strings.Join(fieldSelectors, ",")
	}

	pods, err := internal.GetPods(RestConfig(), podFlags.Namespace, podFlags.Label, podFlags.FieldSelector)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	out.PrintPods(pods)
	return nil
}
