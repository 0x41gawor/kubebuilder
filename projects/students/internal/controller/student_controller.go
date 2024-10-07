/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	schooliov1 "gawor.com/students/api/v1"
)

// StudentReconciler reconciles a Student object
type StudentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=school.io.gawor.com,resources=students,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=school.io.gawor.com,resources=students/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=school.io.gawor.com,resources=students/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Student object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *StudentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Student instance
	var student schooliov1.Student
	if err := r.Get(ctx, req.NamespacedName, &student); err != nil {
		if errors.IsNotFound(err) {
			// Student resource not found, could be deleted, no further action needed
			return ctrl.Result{}, nil
		}
		// Error reading the object, requeue
		return ctrl.Result{}, err
	}

	// Calculate the average grade
	if len(student.Spec.Grades) == 0 {
		logger.Info("No grades found for student", "name", student.Spec.Name)
		return ctrl.Result{}, nil
	}

	var sum int
	for _, grade := range student.Spec.Grades {
		sum += grade
	}
	average := float64(sum) / float64(len(student.Spec.Grades))

	// Determine graduation status
	graduate := average >= 3.0

	// Update status if changed
	if student.Status.Graduate != graduate {
		student.Status.Graduate = graduate
		if err := r.Status().Update(ctx, &student); err != nil {
			logger.Error(err, "Failed to update Student status", "name", student.Spec.Name)
			return ctrl.Result{}, err
		}
	}

	logger.Info("Successfully reconciled Student", "name", student.Spec.Name, "graduate", student.Status.Graduate)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StudentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&schooliov1.Student{}).
		Complete(r)
}
