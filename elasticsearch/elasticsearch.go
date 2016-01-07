package elasticsearch

import (
	"fmt"
	"math"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"

	"github.com/jfklingler/awstagger/context"
)

func Process(ctx context.Context, region string) {
	ctx.Print("  Processing elastic search resources...")
	svc := elasticsearchservice.New(ctx.AwsSession, &aws.Config{Region: aws.String(region)})

	processARNs := func(resourceIds []*string) {
		nextWindow := func(ids []*string) ([]*string, []*string) {
			w := int(math.Min(float64(len(ids)), float64(200)))
			return ids[0:w], ids[w:]
		}

		for thisRound, remaining := nextWindow(resourceIds); len(thisRound) > 0; thisRound, remaining = nextWindow(remaining) {
			updateTags(ctx, *svc, thisRound)
			deleteTags(ctx, *svc, thisRound)
			printTags(ctx, *svc, thisRound)
		}
	}

	processARNs(getArns(svc))
}

func updateTags(ctx context.Context, svc elasticsearchservice.ElasticsearchService, resourceIds []*string) {
	if len(ctx.Tags) <= 0 {
		return
	}

	for _, arn := range resourceIds {
		_, err := svc.AddTags(&elasticsearchservice.AddTagsInput{
			ARN:        arn,
			TagList:    tagArgsToESTags(ctx.Tags),
		})

		kingpin.FatalIfError(err, "Could not update tags for elastic search resources %s", resourceIds)
	}
}

func deleteTags(ctx context.Context, svc elasticsearchservice.ElasticsearchService, resourceIds []*string) {
	if len(ctx.RmTags) <= 0 {
		return
	}

	for _, arn := range resourceIds {
		_, err := svc.RemoveTags(&elasticsearchservice.RemoveTagsInput{
			ARN: 		arn,
			TagKeys:    rmtagArgsToESTags(ctx.RmTags),
		})

		kingpin.FatalIfError(err, "Could not delete tags for elastic search resources %s", resourceIds)
	}
}

func printTags(ctx context.Context, svc elasticsearchservice.ElasticsearchService, resourceIds []*string) {
	if ctx.Verbose {
		for _, arn := range resourceIds {
			ctx.PrintVerbose(fmt.Sprintf("    Processing domain %s", *arn))
			resp, _ := svc.ListTags(&elasticsearchservice.ListTagsInput{
				ARN: arn,
			})

			// Well, this will throw an error if there are no tags instead of just returning an empty list...ugh
			// kingpin.FatalIfError(err, "Could not retrieve tags for elastic search domain %s", arn)

			for _, tag := range resp.TagList {
				ctx.PrintVerbose(fmt.Sprintf("      %s=%s", *tag.Key, *tag.Value))
			}
		}
	}
}

func getArns(svc *elasticsearchservice.ElasticsearchService) []*string {
	nresp, err := svc.ListDomainNames(&elasticsearchservice.ListDomainNamesInput{})

	domainNames := make([]*string, 0)
	for _, dn := range nresp.DomainNames {
		domainNames = append(domainNames, dn.DomainName)
	}

	resp, err := svc.DescribeElasticsearchDomains(&elasticsearchservice.DescribeElasticsearchDomainsInput{
		DomainNames: domainNames,
	})

	kingpin.FatalIfError(err, "Could not retrieve elastic search resources")

	arns := make([]*string, 0)
	for _, esd := range resp.DomainStatusList {
		arns = append(arns, esd.ARN)
	}

	return arns
}

func rmtagArgsToESTags(tags []string) []*string {
	var esTags []*string

	for _, k := range tags {
		kx := k
		esTags = append(esTags, &kx)
	}

	return esTags
}

func tagArgsToESTags(tags map[string]string) []*elasticsearchservice.Tag {
	esTags := make([]*elasticsearchservice.Tag, 0)

	for k, v := range tags {
		kx, vx := k, v
		bob := elasticsearchservice.Tag{
			Key:   &kx,
			Value: &vx,
		}

		esTags = append(esTags, &bob)
	}

	return esTags
}
