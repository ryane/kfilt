package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/ryane/kfilt/pkg/decoder"
	"github.com/ryane/kfilt/pkg/filter"
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
	kind        string
	name        string
	excludeKind string
	excludeName string
	filename    string
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

	rootCmd.Flags().StringVarP(&root.kind, "kind", "k", "", "Only include resources of kind")
	rootCmd.Flags().StringVarP(&root.name, "name", "n", "", "Only include resources with name")
	rootCmd.Flags().StringVarP(&root.excludeKind, "exclude-kind", "K", "", "Exclude resources of kind")
	rootCmd.Flags().StringVarP(&root.excludeName, "exclude-name", "N", "", "Exclude resources with name")
	rootCmd.Flags().StringVarP(&root.filename, "filename", "f", "", "Read manifests from file")

	rootCmd.SetVersionTemplate(`{{.Version}}`)

	return rootCmd
}

func (r *root) run() error {
	var (
		in  []byte
		err error
	)

	// get input
	if r.filename != "" {
		in, err = ioutil.ReadFile(r.filename)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to read file %q", r.filename))
		}
	} else {
		in, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read stdin")
		}
	}

	// decode
	results, err := decoder.New().Decode(in)
	if err != nil {
		return err
	}

	// filter
	filtered := filter.New(
		filter.KindFilter(r.kind),
		filter.NameFilter(r.name),
		filter.ExcludeKindFilter(r.excludeKind),
		filter.ExcludeNameFilter(r.excludeName),
	).Filter(results)

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
