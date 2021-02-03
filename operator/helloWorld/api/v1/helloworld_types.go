/*
Copyright 2020 DesistDaydream.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 编辑这个文件！！！这个就是 Kubebuilder 自动为我们生成的scaffolding(脚手架)
// 该文件主要用来定义 HelloWorld 这个自定义资源应该具有的状态。也就是该资源 yaml 文件中应该写哪些字段。
// 注意：结构体中的 Json 标签是必须的。添加任何新字段都必须具有 json 标签，以便序列化这些字段。
// 注意：每个结构体和其属性的注释，将会被同步到 crd 资源的 description 字段中。
// 可以使用 kustomize build config/crd/ 命令查看当前代码中生成 CRD 的 yaml。

// HelloWorldSpec 定义 HelloWorld 资源的期望状态
// 也就是 HelloWorld 资源 .spec 字段下的字段，比如 pod 资源有 .spec.containers、.spec.volumes 等等字段
type HelloWorldSpec struct {
	// 在这块代码中插入 .spec 字段下的字段，这些就是告诉集群，HelloWorld 资源的期望状态
	// 重要提示: 修改此文件后，一定要运行 make 命令以重新生成代码

	// Foo 是 HelloWorld 资源的 .spec.foo 字段
	Foo string `json:"foo,omitempty"`
}

// HelloWorldStatus 定义 HelloWorld 对象的观察状态
// 也就是 HelloWorld 资源 .status 字段下的字段，主要就是查看当前对象的状态信息。
type HelloWorldStatus struct {
	// 在这块代码中插入 .status 字段下的字段
	// 重要提示: 修改此文件后，一定要运行 make 命令以重新生成代码
}

// +kubebuilder:object:root=true

// HelloWorld 是 helloworld API 的架构
// 就是要想创建 HelloWorld 对象的话，该对象的 yaml 中必须包含的字段，一般都有 apiVersion、kind、metadata、spec、status 这5个。
// metav1.TypeMeta 包含 kind 和 apiVersion 这两个字段下应该定义的内容；metav1.ObjectMeta 包含 metadata 字段下应该定义的内容。
type HelloWorld struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelloWorldSpec   `json:"spec,omitempty"`
	Status HelloWorldStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HelloWorldList 包含 HelloWorld 资源的一个列表
type HelloWorldList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelloWorld `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelloWorld{}, &HelloWorldList{})
}
