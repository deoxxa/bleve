//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package regexp_char_filter

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/couchbaselabs/bleve/analysis"
	"github.com/couchbaselabs/bleve/registry"
)

const Name = "regexp"

type RegexpCharFilter struct {
	r           *regexp.Regexp
	replacement []byte
}

func NewRegexpCharFilter(r *regexp.Regexp, replacement []byte) *RegexpCharFilter {
	return &RegexpCharFilter{
		r:           r,
		replacement: replacement,
	}
}

func (s *RegexpCharFilter) Filter(input []byte) []byte {
	return s.r.ReplaceAllFunc(input, func(in []byte) []byte { return bytes.Repeat(s.replacement, len(in)) })
}

func RegexpCharFilterConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.CharFilter, error) {
	regexpStr, ok := config["regexp"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify regexp")
	}
	r, err := regexp.Compile(regexpStr)
	if err != nil {
		return nil, fmt.Errorf("unable to build regexp char filter: %v", err)
	}
	replaceBytes := []byte(" ")
	replaceStr, ok := config["replace"].(string)
	if ok {
		replaceBytes = []byte(replaceStr)
	}
	return NewRegexpCharFilter(r, replaceBytes), nil
}

func init() {
	registry.RegisterCharFilter(Name, RegexpCharFilterConstructor)
}
