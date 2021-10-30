/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	demov1 "github.com/yashvardhan-kukreja/postgres-writer-operator/api/v1"
	"github.com/yashvardhan-kukreja/postgres-writer-operator/pkg/psql"
)

// PostgresWriterReconciler reconciles a PostgresWriter object
type PostgresWriterReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	PostgresDBClient *psql.PostgresDBClient
}

var (
	finalizers []string = []string{"finalizers.postgreswriters.demo.yash.com/cleanup-row"}
)

//+kubebuilder:rbac:groups=demo.yash.com,resources=postgreswriters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.yash.com,resources=postgreswriters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.yash.com,resources=postgreswriters/finalizers,verbs=update
func (r *PostgresWriterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// capturing the incoming postgres-writer resource
	postgresWriterObject := &demov1.PostgresWriter{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, postgresWriterObject)
	if err != nil {
		// if the resource is not found, then just return (might look redundant as this usually happens in case of Delete events)
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Error occurred while fetching the PostgresWriter resource")
		return ctrl.Result{}, err
	}

	// if the event is not related to delete, just check if the finalizers are rightfully set on the resource
	if postgresWriterObject.GetDeletionTimestamp().IsZero() && !reflect.DeepEqual(finalizers, postgresWriterObject.GetFinalizers()) {
		// set the finalizers of the postgresWriterObject to the rightful ones
		postgresWriterObject.SetFinalizers(finalizers)
		if err := r.Update(ctx, postgresWriterObject); err != nil {
			logger.Error(err, "error occurred while setting the finalizers of the PostgresWriter resource")
			return ctrl.Result{}, err
		}
	}

	// if the metadata.deletionTimestamp is found to be non-zero, this means that the resource is intended and just about to be deleted
	// hence, it's time to clean up the finalizers
	if !postgresWriterObject.GetDeletionTimestamp().IsZero() {
		logger.Info("Deletion detected! Proceeding to cleanup the finalizers...")
		if err := r.cleanupRowFinalizerCallback(ctx, logger, postgresWriterObject); err != nil {
			logger.Error(err, "error occurred while dealing with the cleanup-row finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// sync the resource to the DB - Insert the row if the row doesn't exist in the DB, Update the row if it exists in the DB
	if err := r.synchronizeResourceToDB(logger, *postgresWriterObject); err != nil {
		logger.Error(err, "error occurred while syncing the resource to the DB")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgresWriterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.PostgresWriter{}).
		Complete(r)
}

func (r *PostgresWriterReconciler) synchronizeResourceToDB(logger logr.Logger, postgresWriterObject demov1.PostgresWriter) error {
	id := postgresWriterObject.Namespace + "/" + postgresWriterObject.Name
	table, name, age, country := postgresWriterObject.Spec.Table, postgresWriterObject.Spec.Name, postgresWriterObject.Spec.Age, postgresWriterObject.Spec.Country

	if err := r.PostgresDBClient.Insert(id, table, name, age, country); err != nil {
		return fmt.Errorf("error occurred while writing the row to Postgres DB: %w", err)
	}
	logger.Info(fmt.Sprintf("sychronized the resource %s/%s with the DB successfully", postgresWriterObject.Namespace, postgresWriterObject.Name))
	return nil
}

func (r *PostgresWriterReconciler) cleanupRowFinalizerCallback(ctx context.Context, logger logr.Logger, postgresWriterObject *demov1.PostgresWriter) error {
	// parse the table and the id of the row to delete
	id := postgresWriterObject.Namespace + "/" + postgresWriterObject.Name
	table := postgresWriterObject.Spec.Table

	// delete the row with the above 'id' from the above 'table'
	if err := r.PostgresDBClient.Delete(id, table); err != nil {
		return fmt.Errorf("error occurred while running the DELETE query on Postgres: %w", err)
	}

	// remove the cleanup-row finalizer from the postgresWriterObject
	controllerutil.RemoveFinalizer(postgresWriterObject, "finalizers.postgreswriters.demo.yash.com/cleanup-row")
	if err := r.Update(ctx, postgresWriterObject); err != nil {
		return fmt.Errorf("error occurred while removing the finalizer: %w", err)
	}
	logger.Info("cleaned up the 'finalizers.postgreswriters.demo.yash.com/cleanup-row' finalizer successfully")
	return nil
}
