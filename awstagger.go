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

	"github.com/aws/aws-sdk-go/aws/session"

	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
}

func main() {
	ctx := createContext()

	awsSession := session.New()

	for _, region := range ctx.regions {
		ctx.print(fmt.Sprintf("Processing region %s...", region))

		switch {
		case ctx.doEc2Instances:
			TagInstances(ctx, awsSession, region)
		case ctx.doEc2Amis:
			TagAmis(ctx, awsSession, region)
		}
	}
}
