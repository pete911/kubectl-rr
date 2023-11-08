package cmd

import (
	"context"
	"fmt"
	"github.com/pete911/kubectl-rr/internal/k8s"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
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

	client, err := k8s.NewClient(KubeconfigPath)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := client.GetPods(ctx, podFlags.Namespace, podFlags.Label, podFlags.FieldSelector)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	for _, pod := range pods {
		fmt.Printf("namespace %s pod %s\n", pod.Namespace, pod.Name)
		for _, container := range pod.Containers {
			requests := container.Requests
			limits := container.Limits
			fmt.Printf("  container %s\n", container.Name)
			fmt.Printf("    requests cpu %s memory %s\n", requests.Cpu.String(), requests.Memory.String())
			fmt.Printf("    limites  cpu %s memory %s\n", limits.Cpu.String(), limits.Memory.String())
		}
	}

	return nil
}
