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
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	quiet   = kingpin.Flag("quiet", "Minimal/no output.").Short('q').Bool()
	verbose = kingpin.Flag("verbose", "Verbose output.").Short('v').Bool()

	regions = kingpin.Flag("region", "AWS region to process. (repeatable, default: all standard regions)").Short('r').Strings()
	tags    = kingpin.Flag("tag", "Tag to set/update on all selected resources. (repeatable)").Short('t').PlaceHolder("KEY=VALUE").Strings()
	rmTags  = kingpin.Flag("rm-tag", "Tag key to remove from all selected resources. (repeatable)").PlaceHolder("KEY").Strings()

	doInstances = kingpin.Flag("instances", "Tag EC2 instances. (default: true)").Default("true").Bool()

	allRegions = []string{"us-east-1", "us-west-1", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "sa-east-1"}
)

type context struct {
	quiet   bool
	verbose bool

	regions []string
	tags    map[string]string
	rmTags  []string

	doInstances bool
}

func createContext() context {
	tagMap := make(map[string]string)
	for _, tag := range *tags {
		kv := strings.SplitN(tag, "=", 2)
		tagMap[kv[0]] = kv[1]
	}

	if len(*regions) == 0 {
		regions = &allRegions
	}

	return context{
		quiet:   *quiet,
		verbose: *verbose,

		regions: *regions,
		tags:    tagMap,
		rmTags:  *rmTags,

		doInstances: *doInstances,
	}
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
