/*
Copyright 2016 John Klingler
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
)

type context struct {
	quiet        bool
	verbose      bool
	regions      []string
	tags         []string
	rmTags       []string
	doInstances  bool
}

func (c context) print(line string) {
	if !c.quiet {
		fmt.Println(line)
	}
}

func (c context) printVerbose(line string) {
	if c.verbose && !c.quiet {
		fmt.Println(line)
	}
}
