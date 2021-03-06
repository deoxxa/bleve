# ![bleve](docs/bleve.png) bleve

modern text indexing in go - [blevesearch.com](http://www.blevesearch.com/)

## Features
* Index any go data structure (including JSON)
* Intelligent defaults backed up by powerful configuration
* Supported field types:
    * Text, Numeric, Date
* Supported query types:
    * Term, Phrase, Match, Match Phrase
    * Conjunction, Disjunction, Boolean
    * Numeric Range, Date Range
    * Simple query syntax for human entry
* Search result match highlighting

## Discussion

Discuss usage and development of bleve in the [google group](https://groups.google.com/forum/#!forum/bleve).

## Indexing

		message := struct{
			From: "marty.schoch@gmail.com",
			Body: "bleve indexing is easy",
		}

		mapping := bleve.NewIndexMapping()
		index, _ := bleve.New("example.bleve", mapping)
		index.Index(message)

## Querying

		index, _ := bleve.Open("example.bleve")
		query := bleve.NewSyntaxQuery("bleve")
		searchRequest := bleve.NewSearchRequest(query)
		searchResult, _ := index.Search(searchRequest)


## Status

[![Build Status](https://drone.io/github.com/couchbaselabs/bleve/status.png)](https://drone.io/github.com/couchbaselabs/bleve/latest)
[![Coverage Status](https://coveralls.io/repos/couchbaselabs/bleve/badge.png?branch=master)](https://coveralls.io/r/couchbaselabs/bleve?branch=master)