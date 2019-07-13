package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/ryane/kfilt/pkg/decoder"
	"github.com/ryane/kfilt/pkg/filter"
	"github.com/ryane/kfilt/pkg/input"
	"github.com/ryane/kfilt/pkg/printer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// GitCommit tracks the current git commit
	GitCommit string
	// Version tracks the current version
	Version string
)

type root struct {
	includeKinds []string
	includeNames []string
	excludeKinds []string
	excludeNames []string
	include      []string
	exclude      []string
	filename     string
}

func newRootCommand(args []string) *cobra.Command {
	root := &root{}
	rootCmd := &cobra.Command{
		Use:   "kfilt",
		Short: "kfilt can filter Kubernetes resources",
		Long:  `kfilt can filter Kubernetes resources`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := root.run(); err != nil {
				log.WithError(err).Error()
				os.Exit(1)
			}
		},
		Version: func() string {
			return fmt.Sprintf("%s (%s)\n", Version, GitCommit)
		}(),
	}

	rootCmd.Flags().StringSliceVarP(&root.includeKinds, "kind", "k", []string{}, "Only include resources of kind")
	rootCmd.Flags().StringSliceVarP(&root.includeNames, "name", "n", []string{}, "Only include resources with name")
	rootCmd.Flags().StringSliceVarP(&root.excludeKinds, "exclude-kind", "K", []string{}, "Exclude resources of kind")
	rootCmd.Flags().StringSliceVarP(&root.excludeNames, "exclude-name", "N", []string{}, "Exclude resources with name")
	rootCmd.Flags().StringArrayVarP(&root.include, "include", "i", []string{}, "Include resources matching criteria")
	rootCmd.Flags().StringArrayVarP(&root.exclude, "exclude", "x", []string{}, "Exclude resources matching criteria")
	rootCmd.Flags().StringVarP(&root.filename, "filename", "f", "", "Read manifests from file")

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
		kfilt.AddInclude(filter.Selector{Kind: k})
	}

	for _, n := range r.includeNames {
		kfilt.AddInclude(filter.Selector{Name: n})
	}

	for _, k := range r.excludeKinds {
		kfilt.AddExclude(filter.Selector{Kind: k})
	}

	for _, n := range r.excludeNames {
		kfilt.AddExclude(filter.Selector{Name: n})
	}

	for _, q := range r.include {
		if q != "" {
			s, err := filter.NewSelector(q)
			if err != nil {
				return err
			}
			kfilt.AddInclude(s)
		}
	}

	for _, q := range r.exclude {
		if q != "" {
			s, err := filter.NewSelector(q)
			if err != nil {
				return err
			}
			kfilt.AddExclude(s)
		}
	}

	filtered := kfilt.Filter(results)

	// print
	if err := printer.New().Print(filtered); err != nil {
		return err
	}

	return nil
}

// Execute runs the root command
func Execute(args []string) {
	if err := newRootCommand(args).Execute(); err != nil {
		log.WithError(err).Error()
		os.Exit(2)
	}
}
