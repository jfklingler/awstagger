A simple tool for bulk adding/updating/removing tags from AWS resources.

## Examples
To simply list the tags on all resources in a region:
```bash
awstagger --verbose --region us-east-1
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
