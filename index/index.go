//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package index

import (
	"github.com/couchbaselabs/bleve/document"
)

type Index interface {
	Open() error
	Close()

	Update(doc *document.Document) error
	Delete(id string) error
	Batch(batch Batch) error

	TermFieldReader(term []byte, field string) (TermFieldReader, error)
	DocIdReader(start, end string) (DocIdReader, error)

	FieldReader(field string, startTerm []byte, endTerm []byte) (FieldReader, error)

	DocCount() uint64

	Document(id string) (*document.Document, error)
	DocumentFieldTerms(id string) (FieldTerms, error)

	Fields() ([]string, error)

	SetInternal(key, val []byte) error
	GetInternal(key []byte) ([]byte, error)
	DeleteInternal(key []byte) error

	DumpAll() chan interface{}
	DumpDoc(id string) chan interface{}
	DumpFields() chan interface{}
}

type FieldTerms map[string][]string

type TermFieldVector struct {
	Field string
	Pos   uint64
	Start uint64
	End   uint64
}

type TermFieldDoc struct {
	Term    string
	ID      string
	Freq    uint64
	Norm    float64
	Vectors []*TermFieldVector
}

type TermFieldReader interface {
	Next() (*TermFieldDoc, error)
	Advance(ID string) (*TermFieldDoc, error)
	Count() uint64
	Close()
}

type FieldReader interface {
	Next() (*TermFieldDoc, error)
	Close()
}

type DocIdReader interface {
	Next() (string, error)
	Advance(ID string) (string, error)
	Close()
}

type Batch map[string]*document.Document

func (b Batch) Index(id string, doc *document.Document) {
	b[id] = doc
}

func (b Batch) Delete(id string) {
	b[id] = nil
}
