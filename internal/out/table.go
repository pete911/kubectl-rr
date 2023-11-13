package out

import (
	"fmt"
	"github.com/pete911/kubectl-rr/internal"
	"math"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintPods(pods []internal.Pod) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintln(w, getRow(
		"NAMESPACE", "POD", "CONTAINER",
		"CPU/R", "CPU/L", "CPU", "CPU/Min", "CPU/Max",
		"Memory/R", "Memory/L", "Memory", "Memory/Min", "Memory/Max",
	))
	for _, pod := range pods {
		for _, containers := range getContainerRows(pod) {
			fmt.Fprintln(w, containers)
		}
	}
	fmt.Fprintln(w)
	w.Flush()
}

func getContainerRows(pod internal.Pod) []string {
	var containerRows []string
	for _, container := range pod.InitContainers {
		containerRows = append(containerRows, getContainerRow(pod, container, true))
	}
	for _, container := range pod.Containers {
		containerRows = append(containerRows, getContainerRow(pod, container, false))
	}
	return containerRows
}

func getContainerRow(pod internal.Pod, container internal.Container, init bool) string {
	name := container.Name
	if init {
		name = fmt.Sprintf("%s [init]", name)
	}
	return getRow(
		pod.Namespace, pod.Name, name,
		container.CPU.Request, container.CPU.Limit, formatCPU(container.CPU.Current), formatCPU(container.CPU.Min), formatCPU(container.CPU.Max),
		container.Memory.Request, container.Memory.Limit, formatMemory(container.Memory.Current), formatMemory(container.Memory.Min), formatMemory(container.Memory.Max),
	)
}

func getRow(in ...string) string {
	return strings.Join(in, "\t")
}

func formatCPU(in float64) string {
	if in == 0 {
		return "-"
	}
	if in < 1 {
		return fmt.Sprintf("%.2fm", in*1000)
	}
	return fmt.Sprintf("%.2f", in)
}

func formatMemory(in float64) string {
	if in == 0 {
		return "-"
	}
	return prettyByteSize(in)
}

func prettyByteSize(b float64) string {
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(b) < 1024.0 {
			return fmt.Sprintf("%3.1f%s", b, unit)
		}
		b /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", b)
}
