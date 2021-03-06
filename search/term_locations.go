//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package search

import (
	"sort"
)

type termLocation struct {
	Term  string
	Pos   int
	Start int
	End   int
}

type termLocations []*termLocation

func (t termLocations) Len() int           { return len(t) }
func (t termLocations) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t termLocations) Less(i, j int) bool { return t[i].Start < t[j].Start }

func OrderTermLocations(tlm TermLocationMap) termLocations {
	rv := make(termLocations, 0)
	for term, locations := range tlm {
		for _, location := range locations {
			tl := termLocation{
				Term:  term,
				Pos:   int(location.Pos),
				Start: int(location.Start),
				End:   int(location.End),
			}
			rv = append(rv, &tl)
		}
	}
	sort.Sort(rv)
	return rv
}
