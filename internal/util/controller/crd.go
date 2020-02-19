package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("util")

// RegisterCRD ensures the CRD object is installed into the Kubernetes cluster.
// It will create or update the CRD and it's validation when needed.
func RegisterCRD(ctx context.Context, client apiextclient.Interface,
	crd *apiextv1beta1.CustomResourceDefinition) error {

	scopedLog := log.WithValues("name", crd.Name)
	crdClient := client.ApiextensionsV1beta1().CustomResourceDefinitions()

	clusterCRD, err := crdClient.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		scopedLog.Info("Creating CRD")
		clusterCRD, err = crdClient.Create(crd)
		// This occurs when multiple agents race to create the CRD. Since another has
		// created it, it will also update it, hence the non-error return.
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
	}
	if err != nil {
		return err
	}

	// Update the CRD with the validation schema.
	scopedLog.Info("Updating CRD")
	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		clusterCRD, err = crdClient.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		clusterCRD.Spec = crd.Spec
		if _, err = crdClient.Update(clusterCRD); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		scopedLog.Error(err, "Unable to update CRD")
		return err
	}

	// Wait for the CRD to be available.
	scopedLog.Info("Waiting for CRD to be available")
	if err := WaitForCRD(ctx, client, crd); err != nil {
		if deleteErr := crdClient.Delete(crd.ObjectMeta.Name, nil); deleteErr != nil {
			return fmt.Errorf("unable to delete crd %s: %v "+
				"(deleting CRD due: %v)", crd.Name, deleteErr, err)
		}
		return err
	}

	scopedLog.Info("CRD is installed and up-to-date")
	return nil
}

// RegisterCRDs ensures the CRD objects are installed into the Kubernetes cluster.
// It will create or update the CRDs and theirs validation when needed.
func RegisterCRDs(ctx context.Context, client apiextclient.Interface,
	crds []*apiextv1beta1.CustomResourceDefinition) error {

	for _, crd := range crds {
		if err := RegisterCRD(ctx, client, crd); err != nil {
			return err
		}
	}

	return nil
}

// WaitForCRD waits for the CRD to be available.
func WaitForCRD(ctx context.Context, client apiextclient.Interface,
	crd *apiextv1beta1.CustomResourceDefinition) error {

	scopedLog := log.WithValues("name", crd.Name)
	crdClient := client.ApiextensionsV1beta1().CustomResourceDefinitions()

	return wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err := crdClient.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextv1beta1.Established:
				if cond.Status == apiextv1beta1.ConditionTrue {
					return true, err
				}
			case apiextv1beta1.NamesAccepted:
				if cond.Status == apiextv1beta1.ConditionFalse {
					scopedLog.Error(errors.New(cond.Reason), "Name conflict for CRD")
					return false, err
				}
			}
		}
		return false, err
	})
}

// WaitForCRDs waits for the CRDs to be available.
func WaitForCRDs(ctx context.Context, client apiextclient.Interface,
	crds []*apiextv1beta1.CustomResourceDefinition) error {

	for _, crd := range crds {
		if err := WaitForCRD(ctx, client, crd); err != nil {
			return err
		}
	}

	return nil
}
