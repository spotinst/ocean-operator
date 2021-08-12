// Copyright 2021 NetApp, Inc. All Rights Reserved.

package cli

import (
	"flag"

	"github.com/spotinst/ocean-operator/pkg/log"
	"go.uber.org/zap"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// zapLoggerOptions contains default configuration for the Zap logger.
var zapLoggerOptions = new(ctrlzap.Options)

func init() {
	zapLoggerOptions.BindFlags(flag.CommandLine)
}

func NewZapLogger() log.Logger {
	if zapLoggerOptions.Development {
		zapLoggerOptions.ZapOpts = append(
			zapLoggerOptions.ZapOpts, zap.AddCaller())
	}
	return ctrlzap.New(ctrlzap.UseFlagOptions(zapLoggerOptions))
}
