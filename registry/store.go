//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package registry

import (
	"fmt"

	"github.com/couchbaselabs/bleve/index/store"
)

func RegisterKVStore(name string, constructor KVStoreConstructor) {
	_, exists := stores[name]
	if exists {
		panic(fmt.Errorf("attempted to register duplicate store named '%s'", name))
	}
	stores[name] = constructor
}

type KVStoreConstructor func(config map[string]interface{}) (store.KVStore, error)
type KVStoreRegistry map[string]KVStoreConstructor

func KVStoreConstructorByName(name string) KVStoreConstructor {
	return stores[name]
}
