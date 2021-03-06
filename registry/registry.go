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
	"sort"

	"github.com/couchbaselabs/bleve/analysis"
)

var stores = make(KVStoreRegistry, 0)

// analysis
var charFilters = make(CharFilterRegistry, 0)
var tokenizers = make(TokenizerRegistry, 0)
var tokenMaps = make(TokenMapRegistry, 0)
var tokenFilters = make(TokenFilterRegistry, 0)
var analyzers = make(AnalyzerRegistry, 0)
var dateTimeParsers = make(DateTimeParserRegistry, 0)

type Cache struct {
	CharFilters     CharFilterCache
	Tokenizers      TokenizerCache
	TokenMaps       TokenMapCache
	TokenFilters    TokenFilterCache
	Analyzers       AnalyzerCache
	DateTimeParsers DateTimeParserCache
}

func NewCache() *Cache {
	return &Cache{
		CharFilters:     make(CharFilterCache, 0),
		Tokenizers:      make(TokenizerCache, 0),
		TokenMaps:       make(TokenMapCache, 0),
		TokenFilters:    make(TokenFilterCache, 0),
		Analyzers:       make(AnalyzerCache, 0),
		DateTimeParsers: make(DateTimeParserCache, 0),
	}
}

func (c *Cache) CharFilterNamed(name string) (analysis.CharFilter, error) {
	return c.CharFilters.CharFilterNamed(name, c)
}

func (c *Cache) DefineCharFilter(name string, typ string, config map[string]interface{}) (analysis.CharFilter, error) {
	return c.CharFilters.DefineCharFilter(name, typ, config, c)
}

func (c *Cache) TokenizerNamed(name string) (analysis.Tokenizer, error) {
	return c.Tokenizers.TokenizerNamed(name, c)
}

func (c *Cache) DefineTokenizer(name string, typ string, config map[string]interface{}) (analysis.Tokenizer, error) {
	return c.Tokenizers.DefineTokenizer(name, typ, config, c)
}

func (c *Cache) TokenMapNamed(name string) (analysis.TokenMap, error) {
	return c.TokenMaps.TokenMapNamed(name, c)
}

func (c *Cache) DefineTokenMap(name string, typ string, config map[string]interface{}) (analysis.TokenMap, error) {
	return c.TokenMaps.DefineTokenMap(name, typ, config, c)
}

func (c *Cache) TokenFilterNamed(name string) (analysis.TokenFilter, error) {
	return c.TokenFilters.TokenFilterNamed(name, c)
}

func (c *Cache) DefineTokenFilter(name string, typ string, config map[string]interface{}) (analysis.TokenFilter, error) {
	return c.TokenFilters.DefineTokenFilter(name, typ, config, c)
}

func (c *Cache) AnalyzerNamed(name string) (*analysis.Analyzer, error) {
	return c.Analyzers.AnalyzerNamed(name, c)
}

func (c *Cache) DefineAnalyzer(name string, typ string, config map[string]interface{}) (*analysis.Analyzer, error) {
	return c.Analyzers.DefineAnalyzer(name, typ, config, c)
}

func (c *Cache) DateTimeParserNamed(name string) (analysis.DateTimeParser, error) {
	return c.DateTimeParsers.DateTimeParserNamed(name, c)
}

func (c *Cache) DefineDateTimeParser(name string, typ string, config map[string]interface{}) (analysis.DateTimeParser, error) {
	return c.DateTimeParsers.DefineDateTimeParser(name, typ, config, c)
}

func PrintRegistry() {
	sorted := make(sort.StringSlice, 0, len(charFilters))
	for name, _ := range charFilters {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("Char Filters:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()

	sorted = make(sort.StringSlice, 0, len(tokenizers))
	for name, _ := range tokenizers {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("Tokenizers:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()

	sorted = make(sort.StringSlice, 0, len(tokenMaps))
	for name, _ := range tokenMaps {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("Token Maps:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()

	sorted = make(sort.StringSlice, 0, len(tokenFilters))
	for name, _ := range tokenFilters {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("Token Filters:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()

	sorted = make(sort.StringSlice, 0, len(analyzers))
	for name, _ := range analyzers {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("Analyzers:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()

	sorted = make(sort.StringSlice, 0, len(dateTimeParsers))
	for name, _ := range dateTimeParsers {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("DateTime Parsers:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()

	sorted = make(sort.StringSlice, 0, len(stores))
	for name, _ := range stores {
		sorted = append(sorted, name)
	}
	sorted.Sort()
	fmt.Printf("KV Stores:\n")
	for _, name := range sorted {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println()
}
