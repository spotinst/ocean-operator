// Copyright 2021 NetApp, Inc. All Rights Reserved.

package cli

import (
	"flag"

	"github.com/spotinst/ocean-operator/internal/log"
	"go.uber.org/zap"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// zapLoggerOptions contains default configuration for the Zap logger.
var zapLoggerOptions = &ctrlzap.Options{
	Development: true,
	ZapOpts:     []zap.Option{zap.WithCaller(true)},
}

func init() {
	zapLoggerOptions.BindFlags(flag.CommandLine)
}

func NewZapLogger() log.Logger {
	return ctrlzap.New(ctrlzap.UseFlagOptions(zapLoggerOptions))
}
