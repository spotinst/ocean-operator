// Copyright 2021 NetApp, Inc. All Rights Reserved.

package controllers

import (
	"context"
	"fmt"
	"time"

	oceanv1alpha1 "github.com/spotinst/ocean-operator/api/v1alpha1"
	ctrlutil "github.com/spotinst/ocean-operator/internal/controller"
	"github.com/spotinst/ocean-operator/internal/version"
	"github.com/spotinst/ocean-operator/pkg/installer"
	_ "github.com/spotinst/ocean-operator/pkg/installer/installers"
	"github.com/spotinst/ocean-operator/pkg/log"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// OperatorFinalizerName is the name of the finalizer.
	OperatorFinalizerName = "operator.ocean.spot.io"

	// OperatorVersionAnnotation is the annotation holds the operator version.
	OperatorVersionAnnotation = "operator.ocean.spot.io/version"
)

// OceanComponentReconciler reconciles a OceanComponent object
type OceanComponentReconciler struct {
	Scheme       *runtime.Scheme
	Client       client.Client
	ClientGetter genericclioptions.RESTClientGetter
	Log          log.Logger
	Namespace    string
}

// Helm requires cluster-admin access, but here we'll explicitly mention a few
// resources that the ocean operator accesses directly:
//
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create
// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;uninstall
// +kubebuilder:rbac:groups=ocean.spot.io,resources=components,verbs=get;list;watch;create;update;patch;uninstall
// +kubebuilder:rbac:groups=ocean.spot.io,resources=components/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ocean.spot.io,resources=components/finalizers,verbs=update

// SetupWithManager sets up the controller with the Manager.
func (r *OceanComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oceanv1alpha1.OceanComponent{}).
		Complete(r)
}

type RequestContext struct {
	ctrlutil.RequestContext
	comp      *oceanv1alpha1.OceanComponent
	installer installer.Installer
	log       log.Logger
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.5/pkg/reconcile
func (r *OceanComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rctx := r.newContext(ctx, req)
	rctx.log.Info("reconciling")

	// get component by namespaced name
	rctx.comp = new(oceanv1alpha1.OceanComponent)
	if err := r.Client.Get(ctx, req.NamespacedName, rctx.comp); err != nil {
		if !apierrors.IsNotFound(err) {
			rctx.log.Error(err, "cannot retrieve")
		}
		return ctrlutil.NoRequeue()
	}

	// add finalizer and version annotation
	changed, err := r.setInitialValues(rctx.comp)
	if err != nil {
		return ctrlutil.RequeueError(err)
	}
	if changed {
		if err = r.Client.Update(ctx, rctx.comp); err != nil {
			return ctrlutil.RequeueError(err)
		}
	}

	// initialize new installer
	rctx.installer, err = r.newInstaller(rctx)
	if err != nil {
		rctx.log.Error(err, "cannot reconcile")
		return r.unsupportedType(rctx)
	}

	// reconcile delete
	if ctrlutil.IsBeingDeleted(rctx.comp) {
		resp, err := r.reconcileAbsent(rctx)
		if err != nil {
			return resp, err
		}
		// remove finalizer, but fetch again since it's been patched
		if err = r.Client.Get(ctx, req.NamespacedName, rctx.comp); err != nil {
			if !apierrors.IsNotFound(err) {
				rctx.log.Error(err, "cannot retrieve")
			}
			return ctrlutil.NoRequeue()
		}
		ctrlutil.RemoveFinalizer(rctx.comp, OperatorFinalizerName)
		err = r.Client.Update(ctx, rctx.comp)
		return resp, err
	}

	// reconcile apply
	switch rctx.comp.Spec.State {
	case oceanv1alpha1.OceanComponentStatePresent:
		return r.reconcilePresent(rctx)
	case oceanv1alpha1.OceanComponentStateAbsent:
		return r.reconcileAbsent(rctx)
	default:
		return ctrlutil.RequeueError(fmt.Errorf("unsupported component state: %v", rctx.comp.Spec.State))
	}
}

func (r *OceanComponentReconciler) reconcilePresent(ctx *RequestContext) (ctrl.Result, error) {
	// check whether the component is already installed
	release, err := ctx.installer.Get(ctx.comp.Spec.Name)
	if err != nil {
		if !installer.IsReleaseNotFound(err) {
			return ctrlutil.RequeueError(err)
		} else {
			// component isn't present, install
			return r.install(ctx)
		}
	}

	// component is present, upgrade
	if ctx.installer.IsUpgrade(ctx.comp, release) {
		return r.upgrade(ctx)
	}

	// component is present, and it's not an upgrade
	switch release.Status {
	case installer.ReleaseStatusFailed: // mark as failed, uninstall
		deepCopy := ctx.comp.DeepCopy()
		condition := newCondition(
			oceanv1alpha1.OceanComponentConditionTypeFailure,
			corev1.ConditionTrue,
			installer.ReleaseStatusFailed.String(),
			release.Description)
		changed := setCondition(&(deepCopy.Status), *condition)
		if changed {
			if err = r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
				ctx.log.Error(err, "patch error")
				return ctrlutil.RequeueError(err)
			}
		}
		return r.uninstall(ctx)

	case installer.ReleaseStatusProgressing: // progressing, requeue
		deepCopy := ctx.comp.DeepCopy()
		condition := newCondition(
			oceanv1alpha1.OceanComponentConditionTypeProgressing,
			corev1.ConditionTrue,
			installer.ReleaseStatusProgressing.String(),
			release.Description)
		changed := setCondition(&(deepCopy.Status), *condition)
		if changed {
			if err = r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
				ctx.log.Error(err, "patch error")
				return ctrlutil.RequeueError(err)
			}
		}
		return ctrlutil.RequeueAfterError(15*time.Second, err)

	case installer.ReleaseStatusUninstalled: // well, reinstall it
		deepCopy := ctx.comp.DeepCopy()
		condition := newCondition(
			oceanv1alpha1.OceanComponentConditionTypeAvailable,
			corev1.ConditionFalse,
			installer.ReleaseStatusUninstalled.String(),
			release.Description)
		changed := setCondition(&(deepCopy.Status), *condition)
		if changed { // update component
			if err = r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
				ctx.log.Error(err, "patch error")
				return ctrlutil.RequeueError(err)
			}
			if err = r.Client.Get(ctx, ctx.GetRequest().NamespacedName, ctx.comp); err != nil {
				ctx.log.Error(err, "retrieve error")
				return ctrlutil.RequeueError(err)
			}
		}
		return r.install(ctx)

		// remaining conditions are Deployed and Unknown
		// continue on to component-specific condition
	}

	// check updated conditions
	// note that underlying components may fail without triggering a reconciliation event

	deepCopy := ctx.comp.DeepCopy()
	changed := false

	conditions, err := r.getCurrentConditions(ctx)
	if err != nil {
		ctx.log.Error(err, "cannot get current conditions")
		return ctrlutil.RequeueError(err)
	}
	if conditions != nil {
		for _, condition := range conditions {
			up := setCondition(&(deepCopy.Status), *condition)
			changed = changed || up
		}
	}
	if changed {
		if err = r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}

	condition := getCurrentCondition(deepCopy.Status)
	requeue := true
	if condition.Type == oceanv1alpha1.OceanComponentConditionTypeAvailable &&
		condition.Status == corev1.ConditionTrue {
		requeue = false
	}

	return ctrlutil.Requeue(requeue)
}

func (r *OceanComponentReconciler) reconcileAbsent(ctx *RequestContext) (ctrl.Result, error) {
	_, err := ctx.installer.Get(ctx.comp.Spec.Name)
	if err != nil {
		if installer.IsReleaseNotFound(err) {
			deepCopy := ctx.comp.DeepCopy()
			condition := newCondition(
				oceanv1alpha1.OceanComponentConditionTypeAvailable,
				corev1.ConditionFalse,
				installer.ReleaseStatusUninstalled.String(),
				"Component not present",
			)
			changed := setCondition(&(deepCopy.Status), *condition)
			if changed {
				if err = r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
					ctx.log.Error(err, "patch error")
					return ctrlutil.RequeueError(err)
				}
			}
			return ctrlutil.NoRequeue()
		}
		return ctrlutil.RequeueError(err)
	}

	return r.uninstall(ctx)
}

func (r *OceanComponentReconciler) install(ctx *RequestContext) (ctrl.Result, error) {
	ctx.log.Info("install is required")

	deepCopy := ctx.comp.DeepCopy()
	condition := newCondition(
		oceanv1alpha1.OceanComponentConditionTypeProgressing,
		corev1.ConditionTrue,
		"Installing",
		"Install started",
	)
	changed := setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}

	if err := r.ensureNamespace(ctx, deepCopy.Namespace); err != nil {
		ctx.log.Error(err, "unable to create namespace", "namespace", r.Namespace)
		return ctrlutil.RequeueError(err)
	}

	ephemeralCopy := deepCopy.DeepCopy()
	if err := r.setSensitiveValues(ctx, ephemeralCopy); err != nil {
		return ctrlutil.RequeueError(err)
	}
	_, installErr := ctx.installer.Install(ephemeralCopy)
	if installErr != nil {
		ctx.log.Error(installErr, "installation failed")
		return ctrlutil.RequeueError(installErr)
	}

	condition = newCondition(
		oceanv1alpha1.OceanComponentConditionTypeAvailable,
		corev1.ConditionTrue,
		"Installed",
		"Install finished",
	)
	changed = setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}

	return ctrlutil.RequeueAfter(time.Minute)
}

func (r *OceanComponentReconciler) uninstall(ctx *RequestContext) (ctrl.Result, error) {
	ctx.log.Info("uninstall is required")

	deepCopy := ctx.comp.DeepCopy()
	condition := newCondition(
		oceanv1alpha1.OceanComponentConditionTypeProgressing,
		corev1.ConditionTrue,
		"Uninstalling",
		"Uninstall started",
	)
	changed := setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}
	uninstallErr := ctx.installer.Uninstall(deepCopy)
	if uninstallErr != nil {
		return ctrlutil.RequeueError(uninstallErr)
	}

	condition = newCondition(
		oceanv1alpha1.OceanComponentConditionTypeAvailable,
		corev1.ConditionTrue,
		"Uninstalled",
		"Uninstall finished",
	)
	changed = setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}

	return ctrlutil.RequeueAfter(time.Minute)
}

func (r *OceanComponentReconciler) upgrade(ctx *RequestContext) (ctrl.Result, error) {
	ctx.log.Info("upgrade is required")

	deepCopy := ctx.comp.DeepCopy()
	condition := newCondition(
		oceanv1alpha1.OceanComponentConditionTypeProgressing,
		corev1.ConditionTrue,
		"Upgrading",
		"Upgrade started",
	)
	changed := setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}

	ephemeralCopy := deepCopy.DeepCopy()
	if err := r.setSensitiveValues(ctx, ephemeralCopy); err != nil {
		return ctrlutil.RequeueError(err)
	}
	_, upgradeErr := ctx.installer.Upgrade(ephemeralCopy)
	if upgradeErr != nil {
		return ctrlutil.RequeueError(upgradeErr)
	}

	condition = newCondition(
		oceanv1alpha1.OceanComponentConditionTypeAvailable,
		corev1.ConditionTrue,
		"Upgraded",
		"Upgrade finished",
	)
	changed = setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}

	return ctrlutil.RequeueAfter(time.Minute)
}

func (r *OceanComponentReconciler) ensureNamespace(ctx *RequestContext, namespace string) error {
	if namespace == "" {
		namespace = r.Namespace
	}
	ns := new(corev1.Namespace)
	key := types.NamespacedName{Name: namespace}
	ctx.log.Info("checking existence", "namespace", namespace)
	err := r.Client.Get(ctx, key, ns)
	if apierrors.IsNotFound(err) {
		ns.Name = namespace
		ctx.log.Info("creating", "namespace", namespace)
		return r.Client.Create(ctx, ns)
	}
	return err
}

// getCurrentConditions examines the current environment and returns a condition
// for the component.
func (r *OceanComponentReconciler) getCurrentConditions(ctx *RequestContext) ([]*oceanv1alpha1.OceanComponentCondition, error) {
	objName := ctx.GetRequest().NamespacedName
	switch ctx.comp.Spec.Name {
	case oceanv1alpha1.OceanControllerComponentName, oceanv1alpha1.LegacyOceanControllerComponentName:
		return getOceanControllerConditions(ctx, ctx.log, r.Client, objName)
	case oceanv1alpha1.MetricsServerComponentName:
		return getMetricsServerConditions(ctx, ctx.log, r.Client, objName)
	default:
		// (a) check helm
		// (b) return not installed
		return []*oceanv1alpha1.OceanComponentCondition{
			newCondition(
				oceanv1alpha1.OceanComponentConditionTypeAvailable,
				corev1.ConditionFalse,
				installer.ReleaseStatusUninstalled.String(),
				"Component not installed",
			),
		}, nil
	}
}

func (r *OceanComponentReconciler) unsupportedType(ctx *RequestContext) (ctrl.Result, error) {
	deepCopy := ctx.comp.DeepCopy()
	condition := newCondition(
		oceanv1alpha1.OceanComponentConditionTypeFailure,
		corev1.ConditionTrue,
		installer.ReleaseStatusFailed.String(),
		"Only Helm charts are supported",
	)
	changed := setCondition(&(deepCopy.Status), *condition)
	if changed {
		if err := r.Client.Patch(ctx, deepCopy, client.MergeFrom(ctx.comp)); err != nil {
			ctx.log.Error(err, "patch error")
			return ctrlutil.RequeueError(err)
		}
	}
	return ctrlutil.NoRequeue()
}

func (r *OceanComponentReconciler) setInitialValues(comp *oceanv1alpha1.OceanComponent) (bool, error) {
	changed := false
	if !ctrlutil.IsBeingDeleted(comp) {
		changed = ctrlutil.AddFinalizer(comp, OperatorFinalizerName)
	}
	if comp.Annotations == nil {
		comp.Annotations = make(map[string]string, 1)
	}
	if comp.Annotations[OperatorVersionAnnotation] == "" {
		comp.Annotations[OperatorVersionAnnotation] = version.String()
		changed = true
	}
	return changed, nil
}

func (r *OceanComponentReconciler) setSensitiveValues(ctx *RequestContext, comp *oceanv1alpha1.OceanComponent) error {
	switch comp.Spec.Name {
	case oceanv1alpha1.OceanControllerComponentName, oceanv1alpha1.LegacyOceanControllerComponentName:
		return setOceanControllerSensitiveValues(ctx, ctx.log, r.Client, comp)
	default:
		return nil
	}
}

func (r *OceanComponentReconciler) newContext(ctx context.Context, req ctrl.Request) *RequestContext {
	// generate a new request id
	reqID := ctrlutil.NewRequestId()

	// initialize a new request logger
	reqLog := ctrlutil.NewRequestLog(r.Log, req, reqID)

	// initialize a new base context
	reqCtx := ctrlutil.NewRequestContext(ctx, req, reqID, reqLog)

	// initialize a new request context
	return &RequestContext{
		RequestContext: reqCtx,
		log:            reqLog,
	}
}

func (r *OceanComponentReconciler) newInstaller(ctx *RequestContext) (installer.Installer, error) {
	options := []installer.InstallerOption{
		installer.WithNamespace(r.Namespace),
		installer.WithClientGetter(r.ClientGetter),
		installer.WithLogger(ctx.log),
	}
	switch compType := ctx.comp.Spec.Type; compType {
	case oceanv1alpha1.OceanComponentTypeHelm:
		return installer.GetInstance(string(compType), options...)
	default:
		return nil, fmt.Errorf("unsupported component type: %v", ctx.comp.Spec.Type)
	}
}
