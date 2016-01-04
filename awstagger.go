/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
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
