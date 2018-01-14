/*
Copyright 2017 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	LABEL_APP = "qcloud-app"
)


func int32Ptr(i int32) *int32 { return &i }

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespace := "" //查询所有命名空间下的Pod

	//查找node信息
	nodeClient := clientset.Core().Nodes()
	nodesList, err := nodeClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for i:=0;i<len(nodesList.Items);i++{
		fmt.Printf("node(%d)\n",i)
		fmt.Printf("%#v\n",nodesList.Items[i])
		fmt.Printf("Allocatable: %#v\n",nodesList.Items[i].Status.Allocatable)
	}

	//获取node节点上的Pod
	podClient := clientset.Core().Pods(namespace)
	for i:=0;i<len(nodesList.Items);i++ {
		fmt.Printf("nodeName %s\n",nodesList.Items[i].Name)

		fieldSelector, err := fields.ParseSelector("spec.nodeName=" + nodesList.Items[i].Name + ",status.phase!=" + string(apiv1.PodSucceeded) + ",status.phase!=" + string(apiv1.PodFailed))
		if err != nil {
			panic(err)
		}
		nodeNonTerminatedPodsList, err := podClient.List(metav1.ListOptions{FieldSelector: fieldSelector.String()})
		if err != nil {
			panic(err)
		}

		for j := 0; j < len(nodeNonTerminatedPodsList.Items); j++ {
			fmt.Printf("%s\n", nodeNonTerminatedPodsList.Items[j].Name)
		}
	}

	return
}
