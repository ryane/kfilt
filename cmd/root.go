package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/ryane/kfilt/pkg/decoder"
	"github.com/ryane/kfilt/pkg/filter"
	"github.com/ryane/kfilt/pkg/input"
	"github.com/ryane/kfilt/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	// GitCommit tracks the current git commit
	GitCommit string
	// Version tracks the current version
	Version string
)

type root struct {
	includeKinds         []string
	includeNames         []string
	excludeKinds         []string
	excludeNames         []string
	includeLabelSelector []string
	excludeLabelSelector []string
	include              []string
	exclude              []string
	count                int
	filename             string
}

func newRootCommand(args []string) *cobra.Command {
	root := &root{}
	rootCmd := &cobra.Command{
		Use:   "kfilt",
		Short: "kfilt can filter Kubernetes resources",
		Long:  `kfilt can filter Kubernetes resources`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := root.run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		Version: func() string {
			return fmt.Sprintf("%s (%s)\n", Version, GitCommit)
		}(),
	}

	rootCmd.Flags().StringSliceVarP(&root.includeKinds, "kind", "k", []string{}, "Only include resources of kind")
	rootCmd.Flags().StringSliceVarP(&root.includeNames, "name", "n", []string{}, "Only include resources with name. Wildcards are supported.")
	rootCmd.Flags().StringSliceVarP(&root.excludeKinds, "exclude-kind", "K", []string{}, "Exclude resources of kind")
	rootCmd.Flags().StringSliceVarP(&root.excludeNames, "exclude-name", "N", []string{}, "Exclude resources with name")
	rootCmd.Flags().StringSliceVarP(&root.includeLabelSelector, "labels", "l", []string{}, "Only include resources matching the label selector")
	rootCmd.Flags().StringSliceVarP(&root.excludeLabelSelector, "exclude-labels", "L", []string{}, "Exclude resources matching the label selector")
	rootCmd.Flags().StringArrayVarP(&root.include, "include", "i", []string{}, "Include resources matching criteria")
	rootCmd.Flags().StringArrayVarP(&root.exclude, "exclude", "x", []string{}, "Exclude resources matching criteria")
	rootCmd.Flags().StringVarP(&root.filename, "filename", "f", "", "Read manifests from file or URL")
	rootCmd.Flags().IntVarP(&root.count, "count", "c", 0, "The amount of resources to include")

	rootCmd.SetVersionTemplate(`{{.Version}}`)

	return rootCmd
}

func (r *root) run() error {
	// get input
	in, err := input.Read(r.filename)
	if err != nil {
		if r.filename == "" {
			return errors.Wrap(err, "failed to read stdin")
		}
		return errors.Wrap(err, fmt.Sprintf("failed to read %q", r.filename))
	}
	defer in.Close()

	// decode
	results, err := decoder.New().Decode(in)
	if err != nil {
		return err
	}

	// filter
	kfilt := filter.New()

	for _, k := range r.includeKinds {
		kfilt.AddInclude(filter.Matcher{Kind: k})
	}

	for _, n := range r.includeNames {
		kfilt.AddInclude(filter.Matcher{Name: n})
	}

	for _, l := range r.includeLabelSelector {
		kfilt.AddInclude(filter.Matcher{LabelSelector: l})
	}

	for _, k := range r.excludeKinds {
		kfilt.AddExclude(filter.Matcher{Kind: k})
	}

	for _, n := range r.excludeNames {
		kfilt.AddExclude(filter.Matcher{Name: n})
	}

	for _, l := range r.excludeLabelSelector {
		kfilt.AddExclude(filter.Matcher{LabelSelector: l})
	}

	for _, q := range r.include {
		if q != "" {
			s, err := filter.NewMatcher(q)
			if err != nil {
				return err
			}
			kfilt.AddInclude(s)
		}
	}

	for _, q := range r.exclude {
		if q != "" {
			s, err := filter.NewMatcher(q)
			if err != nil {
				return err
			}
			kfilt.AddExclude(s)
		}
	}

	kfilt.SetCount(r.count)

	filtered, err := kfilt.Filter(results)
	if err != nil {
		return err
	}

	// print
	if err := printer.New().Print(filtered); err != nil {
		return err
	}

	return nil
}

// Execute runs the root command
func Execute(args []string) {
	if err := newRootCommand(args).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
