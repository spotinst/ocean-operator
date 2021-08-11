// Copyright 2021 NetApp, Inc. All Rights Reserved.

package tide

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	ctrlrt "sigs.k8s.io/controller-runtime/pkg/client"
)

func NewConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	cf := genericclioptions.NewConfigFlags(true)
	cf.APIServer = &config.Host
	cf.BearerToken = &config.BearerToken
	cf.CAFile = &config.CAFile
	cf.Namespace = &namespace
	return cf
}

func NewControllerRuntimeClient(config *rest.Config, scheme *runtime.Scheme) (ctrlrt.Client, error) {
	return ctrlrt.New(config, ctrlrt.Options{
		Scheme: scheme,
		Mapper: nil,
	})
}
