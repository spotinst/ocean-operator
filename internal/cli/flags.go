package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spotinst/ocean-operator/pkg/log"
)

// PrintFlags logs the flags in the flag set.
func PrintFlags(flags *pflag.FlagSet, log log.Logger) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.V(4).Info(fmt.Sprintf("FLAG: --%s=%q", flag.Name, flag.Value))
	})
}
