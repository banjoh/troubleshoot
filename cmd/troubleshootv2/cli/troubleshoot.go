package cli

import (
	"os"
	"strings"

	"github.com/replicatedhq/troubleshoot/cmd/util"
	"github.com/replicatedhq/troubleshoot/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

// Load
// Collect
// Analyze
// Redact
// Archive
// Inspect (serve)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "troubleshoot",
		Short: "Troubleshoot commandline interface",
		Long: "A tool for collecting, redacting, analysing support bundles and " +
			"running preflight checks on Kubernetes clusters.\n\n" +
			"For more information, visit https://troubleshoot.sh",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			v := viper.GetViper()
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.BindPFlags(cmd.Flags())

			logger.SetupLogger(v)

			if err := util.StartProfiling(); err != nil {
				klog.Errorf("Failed to start profiling: %v", err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if err := util.StopProfiling(); err != nil {
				klog.Errorf("Failed to stop profiling: %v", err)
			}
		},
	}

	cobra.OnInitialize(initConfig)

	// Subcommands
	cmd.AddCommand(SupporBundleCmd())
	cmd.AddCommand(PreflightCmd())
	cmd.AddCommand(AnalyzeCmd())
	cmd.AddCommand(RedactCmd())
	cmd.AddCommand(InspectCmd())

	// Initialize klog flags
	logger.InitKlogFlags(cmd)

	// CPU and memory profiling flags
	util.AddProfilingFlags(cmd)

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix("TROUBLESHOOT")
	viper.AutomaticEnv()
}
