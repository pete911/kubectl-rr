package cmd

import (
	"fmt"
	"strings"

	"github.com/pete911/kubectl-rr/internal"
	"github.com/pete911/kubectl-rr/internal/k8s"
	"github.com/spf13/cobra"
)

type PodFlags struct {
	Namespace           string
	AllNamespaces       bool
	Label               string
	FieldSelector       string
	PrometheusNamespace string
	PrometheusLabel     string
}

func (f PodFlags) ToConfig(args []string) internal.Config {

	if f.AllNamespaces {
		f.Namespace = ""
	}

	// additional arguments are considered to be pod names, add to field selector flags
	for _, v := range args {
		fieldSelectors := strings.Split(podFlags.FieldSelector, ",")
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.name=%s", v))
		f.FieldSelector = strings.Join(fieldSelectors, ",")
	}

	promConfig := k8s.PrometheusConfig{
		Namespace: f.PrometheusNamespace,
		Labels:    f.PrometheusLabel,
	}

	return internal.Config{
		Namespace:        f.Namespace,
		LabelSelector:    f.Label,
		FieldSelector:    f.FieldSelector,
		PrometheusConfig: promConfig,
	}
}

func InitPodFlags(cmd *cobra.Command, flags *PodFlags) {
	cmd.Flags().StringVarP(
		&flags.Namespace,
		"namespace",
		"n",
		"default",
		"kubernetes namespace",
	)
	cmd.Flags().BoolVarP(
		&flags.AllNamespaces,
		"all-namespaces",
		"A",
		false,
		"all kubernetes namespaces",
	)
	cmd.Flags().StringVarP(
		&flags.Label,
		"label",
		"l",
		"",
		"kubernetes label",
	)
	cmd.Flags().StringVarP(
		&flags.FieldSelector,
		"field-selector",
		"",
		"",
		"kubernetes field selector",
	)
	cmd.Flags().StringVarP(
		&flags.PrometheusNamespace,
		"prometheus-namespace",
		"",
		"",
		"prometheus server namespace",
	)
	cmd.Flags().StringVarP(
		&flags.PrometheusLabel,
		"prometheus-label",
		"",
		"app.kubernetes.io/name=prometheus",
		"prometheus server label selector",
	)
}
