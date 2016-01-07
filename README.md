A simple tool for bulk adding/updating/removing tags from AWS resources.

The ultimate goal is to support all taggable AWS resources, but initially targetting commonly used and "billable"
resources.

Primary benefits over AWS web console or command line tool are:

* No clicking 237 (or 732) checkboxes
* Scriptable
* Simple selectors
* Simple tag specification

## Examples
To simply list the tags on all resources in all regions:
```bash
awstagger --verbose
```

To list the tags on all resources in some regions:
```bash
awstagger --verbose --region us-east-1 --region us-west-2
```

To add tags on all resources in a region:
```bash
awstagger --region us-east-1 --tag key1=value1 --tag key2=value2
```

To remove tags on all resources in a region:
```bash
awstagger --region us-east-1 --rm-tag key3 --rm-tag key4
```

Do it all in one run:
```bash
awstagger --verbose --region us-east-1 --tag key1=value1 --tag key2=value2 --rm-tag key3 --rm-tag key4
```

## Contributing

This is far from a perfect product. Submit issues. Submit pull requests.
