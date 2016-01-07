/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
*/

package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jfklingler/awstagger/context"
	"github.com/jfklingler/awstagger/ec2tags"
	"github.com/jfklingler/awstagger/autoscalingtags"
	"github.com/jfklingler/awstagger/elasticsearch"
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
		autoscalingtags.Process(ctx, region)
		elasticsearch.Process(ctx, region)
	}
}
