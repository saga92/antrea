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

package appliedtogroup

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"

	networkpolicyapis "github.com/vmware-tanzu/antrea/pkg/apis/networkpolicy"
	"github.com/vmware-tanzu/antrea/pkg/apiserver/registry/networkpolicy"
	"github.com/vmware-tanzu/antrea/pkg/apiserver/storage"
	"github.com/vmware-tanzu/antrea/pkg/controller/networkpolicy/store"
	"github.com/vmware-tanzu/antrea/pkg/controller/types"
)

// REST implements rest.Storage for AppliedToGroups.
type REST struct {
	appliedToGroupStore storage.Interface
}

var (
	_ rest.Storage = &REST{}
	_ rest.Watcher = &REST{}
	_ rest.Scoper  = &REST{}
	_ rest.Lister  = &REST{}
	_ rest.Getter  = &REST{}
)

// NewREST returns a REST object that will work against API services.
func NewREST(appliedToGroupStore storage.Interface) *REST {
	return &REST{appliedToGroupStore}
}

func (r *REST) New() runtime.Object {
	return &networkpolicyapis.AppliedToGroup{}
}

func (r *REST) NewList() runtime.Object {
	return &networkpolicyapis.AppliedToGroupList{}
}

func (r *REST) Get(ctx context.Context, name string, options *v1.GetOptions) (runtime.Object, error) {
	addressGroup, exists, err := r.appliedToGroupStore.Get(name)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if !exists {
		return nil, errors.NewNotFound(networkpolicyapis.Resource("appliedtogroup"), name)
	}
	obj := new(networkpolicyapis.AppliedToGroup)
	store.ToAppliedToGroupMsg(addressGroup.(*types.AppliedToGroup), obj, true, nil)
	return obj, nil
}

func (r *REST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	appliedToGroups := r.appliedToGroupStore.List()
	list := new(networkpolicyapis.AppliedToGroupList)
	list.Items = make([]networkpolicyapis.AppliedToGroup, len(appliedToGroups))
	for i := range appliedToGroups {
		store.ToAppliedToGroupMsg(appliedToGroups[i].(*types.AppliedToGroup), &list.Items[i], true, nil)
	}
	return list, nil
}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) Watch(ctx context.Context, options *internalversion.ListOptions) (watch.Interface, error) {
	key, label, field := networkpolicy.GetSelectors(options)
	return r.appliedToGroupStore.Watch(ctx, key, label, field)
}
