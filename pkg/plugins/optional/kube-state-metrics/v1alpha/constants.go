/*
Copyright 2022 The Kubernetes Authors.

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

package v1alpha

// nolint: lll
const MetaDataDescription = `This command will provide kube-state-metrics custom resource config to the project:
  - A yaml file enables ksm to populate state metrics for your CR.
	('kube-state-metrics/cr-config.yaml')

NOTE: This plugin requires:
- kube-state-metrics to be installed in your cluster
- rbac of ksm to access your CRs' status
`
