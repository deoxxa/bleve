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
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

type Person struct {
	Identifier string     `json:"id"`
	Name       string     `json:"name"`
	Age        float64    `json:"age"`
	Title      string     `json:"title"`
	Birthday   time.Time  `json:"birthday"`
	Address    *Address   `json:"address"`
	Hideouts   []*Address `json:"hideouts"`
	Tags       []string   `json:"tags"`
}

func (p *Person) Type() string {
	return "person"
}

var nameMapping = NewDocumentMapping().
	AddFieldMapping(NewFieldMapping("", "text", "en", true, true, true, true))

var tagsMapping = NewDocumentMapping().
	AddFieldMapping(NewFieldMapping("", "text", "standard", true, true, true, false))
var personMapping = NewDocumentMapping().
	AddSubDocumentMapping("name", nameMapping).
	AddSubDocumentMapping("id", NewDocumentDisabledMapping()).
	AddSubDocumentMapping("tags", tagsMapping)

var mapping = NewIndexMapping().
	AddDocumentMapping("person", personMapping)

var people = []*Person{
	&Person{
		Identifier: "a",
		Name:       "marty",
		Age:        19,
		// has no birthday set to test handling of zero time
		Title: "mista",
		Tags:  []string{"gopher", "belieber"},
	},
	&Person{
		Identifier: "b",
		Name:       "steve has a long name",
		Age:        27,
		Birthday:   time.Unix(1000000000, 0),
		Title:      "missess",
	},
	&Person{
		Identifier: "c",
		Name:       "bob walks home",
		Age:        64,
		Birthday:   time.Unix(1400000000, 0),
		Title:      "masta",
	},
	&Person{
		Identifier: "d",
		Name:       "bobbleheaded wings top the phone",
		Age:        72,
		Birthday:   time.Unix(1400000000, 0),
		Title:      "mizz",
	},
}

// FIXME needs more assertions
func TestIndex(t *testing.T) {
	defer os.RemoveAll("testidx")

	index, err := New("testidx", mapping)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// index all the people
	for _, person := range people {
		err = index.Index(person.Identifier, person)
		if err != nil {
			t.Error(err)
		}
	}

	termQuery := NewTermQuery("marti").SetField("name")
	searchRequest := NewSearchRequest(termQuery)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(1) {
		t.Errorf("expected 1 total hit for term query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "a" {
			t.Errorf("expected top hit id 'a', got '%s'", searchResult.Hits[0].ID)
		}
	}

	termQuery = NewTermQuery("noone").SetField("name")
	searchRequest = NewSearchRequest(termQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(0) {
		t.Errorf("expected 0 total hits")
	}

	matchPhraseQuery := NewMatchPhraseQuery("long name")
	searchRequest = NewSearchRequest(matchPhraseQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(1) {
		t.Errorf("expected 1 total hit for phrase query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "b" {
			t.Errorf("expected top hit id 'b', got '%s'", searchResult.Hits[0].ID)
		}
	}

	termQuery = NewTermQuery("walking").SetField("name")
	searchRequest = NewSearchRequest(termQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(0) {
		t.Errorf("expected 0 total hits")
	}

	matchQuery := NewMatchQuery("walking").SetField("name")
	searchRequest = NewSearchRequest(matchQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(1) {
		t.Errorf("expected 1 total hit for match query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "c" {
			t.Errorf("expected top hit id 'c', got '%s'", searchResult.Hits[0].ID)
		}
	}

	prefixQuery := NewPrefixQuery("bobble").SetField("name")
	searchRequest = NewSearchRequest(prefixQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(1) {
		t.Errorf("expected 1 total hit for prefix query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "d" {
			t.Errorf("expected top hit id 'd', got '%s'", searchResult.Hits[0].ID)
		}
	}

	syntaxQuery := NewSyntaxQuery("+name:phone")
	searchRequest = NewSearchRequest(syntaxQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(1) {
		t.Errorf("expected 1 total hit for syntax query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "d" {
			t.Errorf("expected top hit id 'd', got '%s'", searchResult.Hits[0].ID)
		}
	}

	maxAge := 30.0
	numericRangeQuery := NewNumericRangeQuery(nil, &maxAge).SetField("age")
	searchRequest = NewSearchRequest(numericRangeQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(2) {
		t.Errorf("expected 2 total hits for numeric range query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "b" {
			t.Errorf("expected top hit id 'b', got '%s'", searchResult.Hits[0].ID)
		}
		if searchResult.Hits[1].ID != "a" {
			t.Errorf("expected next hit id 'a', got '%s'", searchResult.Hits[1].ID)
		}
	}

	startDate = "2010-01-01"
	dateRangeQuery := NewDateRangeQuery(&startDate, nil).SetField("birthday")
	searchRequest = NewSearchRequest(dateRangeQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(2) {
		t.Errorf("expected 2 total hits for numeric range query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "d" {
			t.Errorf("expected top hit id 'd', got '%s'", searchResult.Hits[0].ID)
		}
		if searchResult.Hits[1].ID != "c" {
			t.Errorf("expected next hit id 'c', got '%s'", searchResult.Hits[1].ID)
		}
	}

	// test that 0 time doesn't get indexed
	endDate = "2010-01-01"
	dateRangeQuery = NewDateRangeQuery(nil, &endDate).SetField("birthday")
	searchRequest = NewSearchRequest(dateRangeQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(1) {
		t.Errorf("expected 1 total hit for numeric range query, got %d", searchResult.Total)
	} else {
		if searchResult.Hits[0].ID != "b" {
			t.Errorf("expected top hit id 'b', got '%s'", searchResult.Hits[0].ID)
		}
	}

	// test behavior of arrays
	// make sure we can successfully find by all elements in array
	termQuery = NewTermQuery("gopher").SetField("tags")
	searchRequest = NewSearchRequest(termQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	} else {
		if searchResult.Total != uint64(1) {
			t.Errorf("expected 1 total hit for term query, got %d", searchResult.Total)
		} else {
			if searchResult.Hits[0].ID != "a" {
				t.Errorf("expected top hit id 'a', got '%s'", searchResult.Hits[0].ID)
			}
		}
	}

	termQuery = NewTermQuery("belieber").SetField("tags")
	searchRequest = NewSearchRequest(termQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	} else {
		if searchResult.Total != uint64(1) {
			t.Errorf("expected 1 total hit for term query, got %d", searchResult.Total)
		} else {
			if searchResult.Hits[0].ID != "a" {
				t.Errorf("expected top hit id 'a', got '%s'", searchResult.Hits[0].ID)
			}
		}
	}

	termQuery = NewTermQuery("notintagsarray").SetField("tags")
	searchRequest = NewSearchRequest(termQuery)
	searchResult, err = index.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	if searchResult.Total != uint64(0) {
		t.Errorf("expected 0 total hits")
	}

	// lookup document a
	// expect to find 2 values for field "tags"
	tagsCount := 0
	doc, err := index.Document("a")
	if err != nil {
		t.Error(err)
	} else {
		for _, f := range doc.Fields {
			if f.Name() == "tags" {
				tagsCount++
			}
		}
	}
	if tagsCount != 2 {
		t.Errorf("expected to find 2 values for tags")
	}
}

func TestIndexCreateNewOverExisting(t *testing.T) {
	defer os.RemoveAll("testidx")

	index, err := New("testidx", NewIndexMapping())
	if err != nil {
		t.Fatal(err)
	}
	index.Close()
	index, err = New("testidx", NewIndexMapping())
	if err != ERROR_INDEX_PATH_EXISTS {
		t.Fatalf("expected error index path exists, got %v", err)
	}
}

func TestIndexOpenNonExisting(t *testing.T) {
	_, err := Open("doesnotexist")
	if err != ERROR_INDEX_PATH_DOES_NOT_EXIST {
		t.Fatalf("expected error index path does not exist, got %v", err)
	}
}

func TestIndexOpenMetaMissingOrCorrupt(t *testing.T) {
	defer os.RemoveAll("testidx")

	index, err := New("testidx", NewIndexMapping())
	if err != nil {
		t.Fatal(err)
	}
	index.Close()

	// now intentionally corrupt the metadata
	ioutil.WriteFile("testidx/index_meta.json", []byte("corrupted"), 0666)

	index, err = Open("testidx")
	if err != ERROR_INDEX_META_CORRUPT {
		t.Fatalf("expected error index metadata corrupted, got %v", err)
	}

	// no intentionally remove the metadata
	os.Remove("testidx/index_meta.json")

	index, err = Open("testidx")
	if err != ERROR_INDEX_META_MISSING {
		t.Fatalf("expected error index metadata missing, got %v", err)
	}
}
