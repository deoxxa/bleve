//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package bleve

import (
	"os"
	"testing"
)

func TestIndexMeta(t *testing.T) {
	var testIndexPath = "doesnotexit.bleve"
	defer os.RemoveAll(testIndexPath)

	// open non-existant meta should error
	_, err := OpenIndexMeta(testIndexPath)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// create meta
	im := &indexMeta{Storage: "leveldb"}
	err = im.Save(testIndexPath)
	if err != nil {
		t.Error(err)
	}
	im = nil

	// open a meta that exists
	im, err = OpenIndexMeta(testIndexPath)
	if err != nil {
		t.Error(err)
	}
	if im.Storage != "leveldb" {
		t.Errorf("expected storage 'leveldb', got '%s'", im.Storage)
	}

	// save a meta that already exists
	err = im.Save(testIndexPath)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
