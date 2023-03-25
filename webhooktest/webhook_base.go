/*
Copyright 2022.

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

package webhooktest

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	testv1beta1 "github.com/gibizer/test-operator/api/v1beta1"
)

const (
	interval = time.Duration(10) * time.Millisecond
	timeout  = interval * 15
)

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
	logger    logr.Logger
)

func CreateNamespace(name string) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	Expect(k8sClient.Create(ctx, ns)).Should(Succeed())
}

func DeleteNamespace(name string) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	Expect(k8sClient.Delete(ctx, ns)).Should(Succeed())
}

func CreateNovaAPI(namespace string, spec map[string]interface{}) types.NamespacedName {
	novaAPIName := uuid.New().String()

	raw := map[string]interface{}{
		"apiVersion": "nova.openstack.org/v1beta1",
		"kind":       "NovaAPI",
		"metadata": map[string]interface{}{
			"name":      novaAPIName,
			"namespace": namespace,
		},
		"spec": spec,
	}
	CreateUnstructured(raw)

	return types.NamespacedName{Name: novaAPIName, Namespace: namespace}
}

func CreateSimple(namespace string) types.NamespacedName {
	name := uuid.New().String()
	raw := map[string]interface{}{
		"apiVersion": "test.test.org/v1beta1",
		"kind":       "Simple",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
	}
	CreateUnstructured(raw)

	logger.Info("Created")
	return types.NamespacedName{Name: name, Namespace: namespace}
}

func DeleteSimple(name types.NamespacedName) {
	Eventually(func(g Gomega) {
		instance := &testv1beta1.Simple{}
		err := k8sClient.Get(ctx, name, instance)
		// if it is already gone that is OK
		if k8s_errors.IsNotFound(err) {
			return
		}
		g.Expect(err).Should(BeNil())

		g.Expect(k8sClient.Delete(ctx, instance)).Should(Succeed())

		err = k8sClient.Get(ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeTrue())
	}, timeout, interval).Should(Succeed())
	logger.Info("Deleted")
}

func GetSimple(name types.NamespacedName) *testv1beta1.Simple {
	instance := &testv1beta1.Simple{}
	Eventually(func(g Gomega) {
		logger.Info("Try get")
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	logger.Info("Get")
	return instance
}

func CreateUnstructured(rawObj map[string]interface{}) {
	logger.Info("Creating", "raw", rawObj)
	unstructuredObj := &unstructured.Unstructured{Object: rawObj}
	_, err := controllerutil.CreateOrPatch(
		ctx, k8sClient, unstructuredObj, func() error { return nil })
	Expect(err).ShouldNot(HaveOccurred())
}
