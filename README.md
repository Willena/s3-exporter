# S3 Exporter

[![dockerhub](https://img.shields.io/docker/v/gillena/s3-exporter?color=blue&label=Docker%20Hub)](https://hub.docker.com/r/gillena/s3-exporter)

A prometheus exporter that expose FS or S3 metrics (file count, volumes, extensions, ...)

## Quick Start

This prometheus exporter expose the following stats

- MaxDepth: Maximum depth of folder tree
- CollectDuration: Time spent reading object and folders
- TotalObjectsSize: Total objects volume in bytes
- TotalObjectsCount: total number of objects found
- LastWalkStart: Date when the stats collection started
- PerPrefixObjectsSizeHistogram: Histogram showing the files size repartition across prefixes
- PerPrefixObjectsSize: Objects volume across prefixes
- PerPrefixObjectsCount: Objects count across prefixes
- PerPrefixPerExtensionObjectCount: Repartition of objects per file extension
- PerPrefixPerExtensionObjectsSize: Total size of objects per extension
- PerPrefixPerContentTypeObjectCount: Repartition of objects per file ContentType
- PerPrefixPerContentTypeObjectsSize: Total size of objects per ContentType

Unlike many prometheus exporters where each http request to scrape metrics triggers collection of them, 
this exporter run using an inner interval that must be set depending on the amount of data that needs to be discovered.

Note: This exporter uses the Minio S3 client, and uses the ListBuckets, ListObjects methods. 

## Options

```
Usage:
  s3-exporter [OPTIONS]

Application Options:
      --type=[s3|fs]                 Walker type [$WALKER_TYPE]
      --interval=                    Define the minimum delay between scrapes.
                                     Set this to a reasonable value to avoid
                                     unnecessary stress on drives (default:
                                     10m) [$SCRAPE_INTERVAL]
      --logLevel=                    Level for logger; available options are:
                                     debug, info, warning, error (default:
                                     debug) [$LOG_LEVEL]

Walkers configuration:
      --walker.maxDepth=             Maximum lookup depth; Will be used to
                                     group paths and results (default: 1)
                                     [$WALKER_MAX_DEPTH]
      --walker.histogram-bins=       Number of bins for histograms (default:
                                     30) [$WALKER_HISTOGRAM_BINS]
      --walker.histogram-start=      Value of first bin in bytes (default:
                                     10_000_000) [$WALKER_HISTOGRAM_START]
      --walker.histogram-factor=     How much do we increase the size of bins
                                     (exponentially) (default: 1.5)
                                     [$WALKER_HISTOGRAM_FACTOR]
      --walker.prefix-filter=        Prefixes or part of prefix to be ignored
                                     [$WALKER_PREFIX_FILTER]
      --walker.custom-labels=        Labels to add for prometheus exporters
                                     [$WALKER_CUSTOM_LABELS]
      --walker.bucket-filter=        Exclude buckets based on name
                                     [$WALKER_BUCKET_FILTER]
      --walker.folder=               Folder to be used for FS walker (default:
                                     /) [$WALKER_FOLDER]

S3 Configuration:
      --walker.s3.endpoint=          URL to the S3 [$WALKER_S3_ENDPOINT]
      --walker.s3.bucket=            S3 bucket [$WALKER_S3_BUCKET]
      --walker.s3.access-key=        S3 Storage Access Key
                                     [$WALKER_S3_ACCESS_KEY]
      --walker.s3.secret-key=        S3 Storage Secret Key
                                     [$WALKER_S3_SECRET_KEY]
      --walker.s3.region=            S3 Storage Region (default: us-west)
                                     [$WALKER_S3_REGION]
      --walker.s3.bucket-path-style  Bucket type [$WALKER_S3_BUCKET_PATH_STYLE]

HTTP Server configuration:
      --http.port=                   HTTP(s) server port (default: 6535) [$PORT]
      --http.addr=                   HTTP(s) listen address [$ADDR]
      --http.keyFile=                Required along with certFile to enable
                                     HTTPS [$KEY_FILE]
      --http.certFile=               Required along with keyFile to enable
                                     HTTPS [$CERT_FILE]

Help Options:
  -h, --help                         Show this help message
```

## License

Copyright 2021 Guillaume VILLENA (Willena)

Licensed under the Apache License, Version 2.0 (the "License")
