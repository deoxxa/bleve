//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package analysis

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

type TokenMap map[string]bool

func NewTokenMap() TokenMap {
	return make(TokenMap, 0)
}

func (s TokenMap) LoadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return s.LoadBytes(data)
}

func (t TokenMap) LoadBytes(data []byte) error {
	bytesReader := bytes.NewReader(data)
	bufioReader := bufio.NewReader(bytesReader)
	line, err := bufioReader.ReadString('\n')
	for err == nil {
		t.LoadLine(line)
		line, err = bufioReader.ReadString('\n')
	}
	// if the err was EOF still need to process last value
	if err == io.EOF {
		t.LoadLine(line)
		return nil
	}
	return err
}

func (t TokenMap) LoadLine(line string) error {
	// find the start of comment, if any
	startComment := strings.IndexAny(line, "#|")
	if startComment >= 0 {
		line = line[:startComment]
	}

	tokens := strings.Fields(line)
	for _, token := range tokens {
		t.AddToken(token)
	}
	return nil
}

func (t TokenMap) AddToken(token string) {
	t[token] = true
}
