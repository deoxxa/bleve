//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package document

import (
	"fmt"

	"github.com/couchbaselabs/bleve/analysis"
)

const DEFAULT_TEXT_INDEXING_OPTIONS = INDEX_FIELD

type TextField struct {
	name           string
	arrayPositions []uint64
	options        IndexingOptions
	analyzer       *analysis.Analyzer
	value          []byte
}

func (t *TextField) Name() string {
	return t.name
}

func (t *TextField) ArrayPositions() []uint64 {
	return t.arrayPositions
}

func (t *TextField) Options() IndexingOptions {
	return t.options
}

func (t *TextField) Analyze() (int, analysis.TokenFrequencies) {
	var tokens analysis.TokenStream
	if t.analyzer != nil {
		tokens = t.analyzer.Analyze(t.Value())
	} else {
		tokens = analysis.TokenStream{
			&analysis.Token{
				Start:    0,
				End:      len(t.value),
				Term:     t.value,
				Position: 1,
				Type:     analysis.AlphaNumeric,
			},
		}
	}
	fieldLength := len(tokens) // number of tokens in this doc field
	tokenFreqs := analysis.TokenFrequency(tokens)
	return fieldLength, tokenFreqs
}

func (t *TextField) Value() []byte {
	return t.value
}

func (t *TextField) GoString() string {
	return fmt.Sprintf("&document.TextField{Name:%s, Options: %s, Analyzer: %s, Value: %s}", t.name, t.options, t.analyzer, t.value)
}

func NewTextField(name string, arrayPositions []uint64, value []byte) *TextField {
	return NewTextFieldWithIndexingOptions(name, arrayPositions, value, DEFAULT_TEXT_INDEXING_OPTIONS)
}

func NewTextFieldWithIndexingOptions(name string, arrayPositions []uint64, value []byte, options IndexingOptions) *TextField {
	return &TextField{
		name:           name,
		arrayPositions: arrayPositions,
		options:        options,
		value:          value,
	}
}

func NewTextFieldWithAnalyzer(name string, arrayPositions []uint64, value []byte, analyzer *analysis.Analyzer) *TextField {
	return &TextField{
		name:           name,
		arrayPositions: arrayPositions,
		options:        DEFAULT_TEXT_INDEXING_OPTIONS,
		analyzer:       analyzer,
		value:          value,
	}
}

func NewTextFieldCustom(name string, arrayPositions []uint64, value []byte, options IndexingOptions, analyzer *analysis.Analyzer) *TextField {
	return &TextField{
		name:           name,
		arrayPositions: arrayPositions,
		options:        options,
		analyzer:       analyzer,
		value:          value,
	}
}