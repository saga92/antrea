// Copyright 2019 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"fmt"
	"testing"

	"github.com/vmware-tanzu/antrea/pkg/agent/util"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteNetworkpolicy(data *TestData, policy *networkingv1.NetworkPolicy) error {
	if err := data.clientset.NetworkingV1().NetworkPolicies(policy.Namespace).Delete(policy.Name, nil); err != nil {
		return fmt.Errorf("unable to cleanup policy %v: %v", policy.Name, err)
	}
	return nil
}

func TestDenyAllNetworkpolicy(t *testing.T) {
	data, err := setupTest(t)
	if err != nil {
		t.Fatalf("Error when setting up test: %v", err)
	}
	defer teardownTest(t, data)

	ns := "ns-e2e-network-policy"

	policy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deny-all",
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress:     []networkingv1.NetworkPolicyIngressRule{},
		},
	}

	policy, err = data.clientset.NetworkingV1().NetworkPolicies(ns).Create(policy)
	defer func() {
		err = deleteNetworkpolicy(data, policy)
		if err != nil {
			t.Fatalf("Error when deleting network policy: %v", err)
		}
	}()

	workerNode := workerNodeName(1)
	podName1 := "pod1"
	if err := data.createBusyboxPodOnNode(podName1, workerNode); err != nil {
		t.Fatalf("Error when creating busybox test Pod: %v", err)
	}
	defer deletePodWrapper(t, data, podName1)

	podName2 := "pod2"
	if err := data.createBusyboxPodOnNode(podName2, workerNode); err != nil {
		t.Fatalf("Error when creating busybox test Pod: %v", err)
	}
	defer deletePodWrapper(t, data, podName2)

	// if err = data.runPingCommandFromTestPod(podName1, podName2, 10); err == nil {
	// 	t.Fatalf("Two pods should not be connected: %v", err)
	// }
}

// func TestBasicNetworkpolicy(t *testing.T) {
// 	data, err := setupTest(t)
// 	if err != nil {
// 		t.Fatalf("Error when setting up test: %v", err)
// 	}
// 	defer teardownTest(t, data)

// 	pod1Name := "pod1"
// 	namespace1 := "namespace1"

// 	policy := &networkingv1.NetworkPolicy{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "allow-client-a-via-pod-selector-with-match-expressions",
// 		},
// 		Spec: networkingv1.NetworkPolicySpec{
// 			PodSelector: metav1.LabelSelector{
// 				MatchLabels: map[string]string{
// 					"pod-name": pod1Name,
// 				},
// 			},
// 			Ingress: []networkingv1.NetworkPolicyIngressRule{{
// 				From: []networkingv1.NetworkPolicyPeer{{
// 					PodSelector: &metav1.LabelSelector{
// 						MatchExpressions: []metav1.LabelSelectorRequirement{{
// 							Key:      "pod-name",
// 							Operator: metav1.LabelSelectorOpIn,
// 							Values:   []string{"client-a"},
// 						}},
// 					},
// 				}},
// 			}},
// 		},
// 	}

// 	policy, err = data.clientset.NetworkingV1().NetworkPolicies(namespace1).Create(policy)

// }