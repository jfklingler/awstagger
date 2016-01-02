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

	"github.com/jfklingler/awstagger/context"
	"github.com/jfklingler/awstagger/ec2tags"

	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
}

func main() {
	ctx := context.New()

	for _, region := range ctx.Regions {
		ctx.Print(fmt.Sprintf("Processing region %s...", region))

		ec2tags.Process(ctx, region)
	}
}
