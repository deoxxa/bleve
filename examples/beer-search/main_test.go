package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/couchbaselabs/bleve"
)

func TestBeerSearchAll(t *testing.T) {
	defer os.RemoveAll("beer-search-test.bleve")

	mapping := buildIndexMapping()
	index, err := bleve.New("beer-search-test.bleve", mapping)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	for jf := range walkDirectory("../../samples/beer-sample/", t) {
		docId := jf.filename[0:strings.LastIndex(jf.filename, ".")]
		err = index.Index(docId, jf.contents)
		if err != nil {
			t.Error(err)
		}
	}

	expectedCount := uint64(7303)
	actualCount := index.DocCount()
	if actualCount != expectedCount {
		t.Errorf("expected %d documents, got %d", expectedCount, actualCount)
	}

	// run a term search
	termQuery := bleve.NewTermQuery("shock").SetField("name")
	termSearchRequest := bleve.NewSearchRequest(termQuery)
	termSearchResult, err := index.Search(termSearchRequest)
	if err != nil {
		t.Error(err)
	}
	expectedResultCount := uint64(1)
	if termSearchResult.Total != expectedResultCount {
		t.Errorf("expected %d hits, got %d", expectedResultCount, termSearchResult.Total)
	} else {
		expectedResultId := "anheuser_busch-shock_top"
		if termSearchResult.Hits[0].ID != expectedResultId {
			t.Errorf("expected top hit ID: %s, got %s", expectedResultId, termSearchResult.Hits[0].ID)
		}
	}

	// run a match phrase search
	matchPhraseQuery := bleve.NewMatchPhraseQuery("spicy mexican food")
	matchPhraseSearchRequest := bleve.NewSearchRequest(matchPhraseQuery)
	matchPhraseSearchResult, err := index.Search(matchPhraseSearchRequest)
	if err != nil {
		t.Error(err)
	}
	expectedResultCount = uint64(1)
	if matchPhraseSearchResult.Total != expectedResultCount {
		t.Errorf("expected %d hits, got %d", expectedResultCount, matchPhraseSearchResult.Total)
	} else {
		expectedResultId := "great_divide_brewing-wild_raspberry_ale"
		if matchPhraseSearchResult.Hits[0].ID != expectedResultId {
			t.Errorf("expected top hit ID: %s, got %s", expectedResultId, matchPhraseSearchResult.Hits[0].ID)
		}
	}

	// run a syntax query
	syntaxQuery := bleve.NewSyntaxQuery("+name:light +description:water -description:barley")
	syntaxSearchRequest := bleve.NewSearchRequest(syntaxQuery)
	syntaxSearchResult, err := index.Search(syntaxSearchRequest)
	if err != nil {
		t.Error(err)
	}
	expectedResultCount = uint64(1)
	if syntaxSearchResult.Total != expectedResultCount {
		t.Errorf("expected %d hits, got %d", expectedResultCount, syntaxSearchResult.Total)
	} else {
		expectedResultId := "iron_city_brewing_co-ic_light"
		if syntaxSearchResult.Hits[0].ID != expectedResultId {
			t.Errorf("expected top hit ID: %s, got %s", expectedResultId, syntaxSearchResult.Hits[0].ID)
		}
	}

	// run a numeric range search
	queryMin := 50.0
	numericRangeQuery := bleve.NewNumericRangeQuery(&queryMin, nil).SetField("abv")
	numericSearchRequest := bleve.NewSearchRequest(numericRangeQuery)
	numericSearchResult, err := index.Search(numericSearchRequest)
	if err != nil {
		t.Error(err)
	}
	expectedResultCount = uint64(1)
	if numericSearchResult.Total != expectedResultCount {
		t.Errorf("expected %d hits, got %d", expectedResultCount, numericSearchResult.Total)
	} else {
		expectedResultId := "woodforde_s_norfolk_ales-norfolk_nog_old_dark_ale"
		if numericSearchResult.Hits[0].ID != expectedResultId {
			t.Errorf("expected top hit ID: %s, got %s", expectedResultId, numericSearchResult.Hits[0].ID)
		}
	}

	// run a date range search
	queryStartDate := "2011-10-04"
	dateRangeQuery := bleve.NewDateRangeQuery(&queryStartDate, nil).SetField("updated")
	dateSearchRequest := bleve.NewSearchRequest(dateRangeQuery)
	dateSearchResult, err := index.Search(dateSearchRequest)
	if err != nil {
		t.Error(err)
	}
	expectedResultCount = uint64(2)
	if dateSearchResult.Total != expectedResultCount {
		t.Errorf("expected %d hits, got %d", expectedResultCount, dateSearchResult.Total)
	} else {
		expectedResultId := "brasserie_du_bouffay-ambr"
		if dateSearchResult.Hits[0].ID != expectedResultId {
			t.Errorf("expected top hit ID: %s, got %s", expectedResultId, dateSearchResult.Hits[0].ID)
		}
	}

	// run a prefix search
	prefixQuery := bleve.NewPrefixQuery("adir").SetField("name")
	prefixSearchRequest := bleve.NewSearchRequest(prefixQuery)
	prefixSearchResult, err := index.Search(prefixSearchRequest)
	if err != nil {
		t.Error(err)
	}
	expectedResultCount = uint64(1)
	if prefixSearchResult.Total != expectedResultCount {
		t.Errorf("expected %d hits, got %d", expectedResultCount, prefixSearchResult.Total)
	} else {
		expectedResultId := "f_x_matt_brewing-saranac_adirondack_lager"
		if prefixSearchResult.Hits[0].ID != expectedResultId {
			t.Errorf("expected top hit ID: %s, got %s", expectedResultId, prefixSearchResult.Hits[0].ID)
		}
	}
}

type jsonFile struct {
	filename string
	contents []byte
}

func walkDirectory(dir string, t *testing.T) chan jsonFile {
	rv := make(chan jsonFile)
	go func() {
		defer close(rv)

		// open the directory
		dirEntries, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		// walk the directory entries
		for _, dirEntry := range dirEntries {
			// read the bytes
			jsonBytes, err := ioutil.ReadFile(dir + "/" + dirEntry.Name())
			if err != nil {
				t.Fatal(err)
			}

			rv <- jsonFile{
				filename: dirEntry.Name(),
				contents: jsonBytes,
			}
		}
	}()
	return rv
}
