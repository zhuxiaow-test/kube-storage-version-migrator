/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"
	"reflect"

	migration_v1alpha1 "github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/apis/migration/v1alpha1"
	migrationclient "github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/clientset"
	migrationinformer "github.com/kubernetes-sigs/kube-storage-version-migrator/pkg/informer/migration/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	StatusIndex     = "Status"
	StatusRunning   = "Running"
	StatusPending   = "Pending"
	StatusCompleted = "Completed"
)

// migrationStatusIndexFunc categorize StorageVersionMigrations based on their conditions.
func migrationStatusIndexFunc(obj interface{}) ([]string, error) {
	m, ok := obj.(*migration_v1alpha1.StorageVersionMigration)
	if !ok {
		return []string{}, fmt.Errorf("expected StroageVersionMigration, got %#v", reflect.TypeOf(obj))
	}
	if HasCondition(m, migration_v1alpha1.MigrationSucceeded) || HasCondition(m, migration_v1alpha1.MigrationFailed) {
		return []string{StatusCompleted}, nil
	}
	if HasCondition(m, migration_v1alpha1.MigrationRunning) {
		return []string{StatusRunning}, nil
	}
	return []string{StatusPending}, nil
}

func NewStatusIndexedInformer(c migrationclient.Interface) cache.SharedIndexInformer {
	return migrationinformer.NewStorageVersionMigrationInformer(c, metav1.NamespaceAll, 0, cache.Indexers{StatusIndex: migrationStatusIndexFunc})
}
