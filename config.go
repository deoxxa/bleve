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
	"expvar"
	"time"

	"github.com/couchbaselabs/bleve/search"

	// token filters
	_ "github.com/couchbaselabs/bleve/analysis/token_filters/lower_case_filter"
	_ "github.com/couchbaselabs/bleve/analysis/token_filters/ngram_filter"

	// tokenizers
	_ "github.com/couchbaselabs/bleve/analysis/tokenizers/single_token"

	// date time parsers
	_ "github.com/couchbaselabs/bleve/analysis/datetime_parsers/datetime_optional"
	_ "github.com/couchbaselabs/bleve/analysis/datetime_parsers/flexible_go"

	// kv stores
	_ "github.com/couchbaselabs/bleve/index/store/boltdb"
	_ "github.com/couchbaselabs/bleve/index/store/inmem"
)

var bleveExpVar = expvar.NewMap("bleve")

type HighlightConfig struct {
	Highlighters map[string]search.Highlighter
}

type Configuration struct {
	Highlight           *HighlightConfig
	DefaultHighlighter  *string
	ByteArrayConverters map[string]ByteArrayConverter
	DefaultKVStore      string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Highlight: &HighlightConfig{
			Highlighters: make(map[string]search.Highlighter),
		},
		ByteArrayConverters: make(map[string]ByteArrayConverter),
	}
}

var Config *Configuration

func init() {
	bootStart := time.Now()

	// build the default configuration
	Config = NewConfiguration()

	// register byte array converters
	Config.ByteArrayConverters["string"] = NewStringByteArrayConverter()
	Config.ByteArrayConverters["json"] = NewJSONByteArrayConverter()
	Config.ByteArrayConverters["ignore"] = NewIgnoreByteArrayConverter()

	// register ansi highlighter
	Config.Highlight.Highlighters["ansi"] = search.NewSimpleHighlighter()

	// register html highlighter
	htmlFormatter := search.NewHTMLFragmentFormatterCustom(`<span class="highlight">`, `</span>`)
	htmlHighlighter := search.NewSimpleHighlighter()
	htmlHighlighter.SetFragmentFormatter(htmlFormatter)
	Config.Highlight.Highlighters["html"] = htmlHighlighter

	// set the default highlighter
	htmlHighlighterName := "html"
	Config.DefaultHighlighter = &htmlHighlighterName

	// default kv store
	Config.DefaultKVStore = "boltdb"

	bootDuration := time.Since(bootStart)
	bleveExpVar.Add("bootDuration", int64(bootDuration))
}
