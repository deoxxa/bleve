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
	"github.com/couchbaselabs/bleve/search"
)

type BooleanQuery struct {
	Must     *ConjunctionQuery `json:"must,omitempty"`
	Should   *DisjunctionQuery `json:"should,omitempty"`
	MustNot  *DisjunctionQuery `json:"must_not,omitempty"`
	BoostVal float64           `json:"boost,omitempty"`
}

func NewBooleanQuery(must *ConjunctionQuery, should *DisjunctionQuery, mustNot *DisjunctionQuery) *BooleanQuery {
	return &BooleanQuery{
		Must:     must,
		Should:   should,
		MustNot:  mustNot,
		BoostVal: 1.0,
	}
}

func (q *BooleanQuery) Boost() float64 {
	return q.BoostVal
}

func (q *BooleanQuery) SetBoost(b float64) *BooleanQuery {
	q.BoostVal = b
	return q
}

func (q *BooleanQuery) Searcher(i *indexImpl, explain bool) (search.Searcher, error) {

	var err error
	var mustSearcher *search.ConjunctionSearcher
	if q.Must != nil {
		mustSearcher, err = q.Must.Searcher(i, explain)
		if err != nil {
			return nil, err
		}
	}

	var shouldSearcher *search.DisjunctionSearcher
	if q.Should != nil {
		shouldSearcher, err = q.Should.Searcher(i, explain)
		if err != nil {
			return nil, err
		}
	}

	var mustNotSearcher *search.DisjunctionSearcher
	if q.MustNot != nil {
		mustNotSearcher, err = q.MustNot.Searcher(i, explain)
		if err != nil {
			return nil, err
		}
	}

	return search.NewBooleanSearcher(i.i, mustSearcher, shouldSearcher, mustNotSearcher, explain)
}

func (q *BooleanQuery) Validate() error {
	if q.Must != nil {
		err := q.Must.Validate()
		if err != nil {
			return err
		}
	}
	if q.Should != nil {
		err := q.Should.Validate()
		if err != nil {
			return err
		}
	}
	if q.MustNot != nil {
		err := q.MustNot.Validate()
		if err != nil {
			return err
		}
	}
	if q.Must == nil && q.Should == nil {
		return ERROR_BOOLEAN_QUERY_NEEDS_MUST_OR_SHOULD
	}
	if q.Must != nil && len(q.Must.Conjuncts) == 0 && q.Should != nil && len(q.Should.Disjuncts) == 0 {
		return ERROR_BOOLEAN_QUERY_NEEDS_MUST_OR_SHOULD
	}
	return nil
}
