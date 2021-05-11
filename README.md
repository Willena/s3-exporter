# S3 Exporter

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
./s3_exporter [OPTIONS]

Application Options:
--type:[s3|fs]                  Walker type [env: WALKER_TYPE]
--interval:                     Define the minimum delay between scrapes. Set this to a reasonable value to avoid unnecessary stress on drives (default: 10m) [SCRAPE_INTERVAL]
--logLevel:                     Level for logger; available options are: debug, info, warning, error (default: debug) [LOG_LEVEL]

Walkers configuration:
--walker.maxDepth:              Maximum lookup depth; Will be used to group paths and results (default: 1) [env: WALKER_MAX_DEPTH]
--walker.histogram-bins:        Number of bins for histograms (default: 30) [env: WALKER_HISTOGRAM_BINS]
--walker.histogram-start:       Value of first bin in bytes (default: 10_000_000) [env: WALKER_HISTOGRAM_START]
--walker.histogram-factor:      How much do we increase the size of bins (exponentially) (default: 1.5) [env: WALKER_HISTOGRAM_FACTOR]
--walker.prefix-filter:         Prefixes or part of prefix to be ignored [env: WALKER_PREFIX_FILTER]
--walker.custom-labels:         Labels to add for prometheus exporters [env: WALKER_CUSTOM_LABELS]
--walker.bucket-filter:         Exclude buckets based on name [env: WALKER_BUCKET_FILTER]
--walker.folder:                Folder to be used for FS walker (default: /) [env: WALKER_FOLDER]

S3 Configuration:
--walker.S3.endpoint:           URL to the S3 [env: WALKER_S3_ENDPOINT]
--walker.S3.bucket:             S3 bucket [env: WALKER_S3_BUCKET]
--walker.S3.class:              S3 Storage Class (default: STANDARD) [env: WALKER_S3_CLASS]
--walker.S3.access-key:         S3 Storage Access Key [env: WALKER_S3_ACCESS_KEY]
--walker.S3.token:              S3 Access token [env: WALKER_S3_TOKEN]
--walker.S3.secret-key:         S3 Storage Secret Key [env: WALKER_S3_SECRET_KEY]
--walker.S3.region:             S3 Storage Region (default: us-west) [env: WALKER_S3_REGION]
--walker.S3.bucket-path-style   Bucket type [env: WALKER_S3_BUCKET_PATH_STYLE]

HTTP Server configuration:
--http.port:                    HTTP(s) server port (default: 6535) [env: PORT]
--http.addr:                    HTTP(s) listen address [env: ADDR]
--http.keyFile:                 Required along with certFile to enable HTTPS [env: KEY_FILE]
--http.certFile:                Required along with keyFile to enable HTTPS [env: CERT_FILE]

Help Options:
-?                                 Show this help message
-h, --help                          Show this help message
```

## License

Copyright 2021 Guillaume VILLENA (Willena)

Licensed under the Apache License, Version 2.0 (the "License")