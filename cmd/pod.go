package cmd

import (
	"fmt"
	"github.com/pete911/kubectl-rr/internal"
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

	for _, pod := range pods {
		fmt.Printf("namespace %s pod %s\n", pod.Namespace, pod.Name)
		for _, container := range pod.Containers {
			fmt.Printf("  container %s\n", container.Name)
			fmt.Printf("    cpu requests %s limits %s current: %s min: %s max: %s\n",
				container.CPU.Request, container.CPU.Limit, container.CPU.Current, container.CPU.Min, container.CPU.Max)
		}
	}

	return nil
}
