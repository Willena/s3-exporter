package walker

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/willena/s3-exporter/utils"
	"net/url"
	"regexp"
)

type S3WalkerConfig struct {
	S3Configuration `group:"S3 Configuration" namespace:"s3" env-namespace:"S3"`
	BucketFilters   []string `long:"bucket-filter" env:"BUCKET_FILTER" description:"Exclude buckets based on name"`
}

type S3Configuration struct {
	Endpoint        string `long:"endpoint" description:"URL to the S3" required:"false" env:"ENDPOINT"`
	Bucket          string `long:"bucket" description:"S3 bucket" required:"false" env:"BUCKET"`
	AccessKey       string `long:"access-key" description:"S3 Storage Access Key" required:"false" env:"ACCESS_KEY"`
	SecretKey       string `long:"secret-key" description:"S3 Storage Secret Key" required:"false" env:"SECRET_KEY"`
	Region          string `long:"region" description:"S3 Storage Region" required:"false" env:"REGION" default:"us-west"`
	BucketPathStyle bool   `long:"bucket-path-style" description:"Bucket type" required:"false" env:"BUCKET_PATH_STYLE"`
}

type S3Walker struct {
	baseWalker
	config         *S3WalkerConfig
	client         *minio.Client
	bucketPatterns []*regexp.Regexp
}

func (s *S3Walker) Init(config Config, labels map[string]string, _ []string) error {
	err := s.ValidateConfig(config)
	if err != nil {
		return err
	}
	s.config = config.S3WalkerConfig
	s.client = s.createClient()

	s.bucketPatterns = utils.BuildPatternsFromStrings(s.config.BucketFilters)

	return s.baseWalker.Init(config,
		utils.MergeMapsRight(map[string]string{
			"type":       "s3Walker",
			"s3Endpoint": s.config.Endpoint,
		}, labels), []string{"bucket", "storageClass"})
}

func (s *S3Walker) createClient() *minio.Client {
	// Initialize minio client object.
	uri, err := url.ParseRequestURI(s.config.Endpoint)
	if err != nil {
		log.Fatalln("Could not read S3 url")
	}

	bucketType := minio.BucketLookupDNS
	if s.config.BucketPathStyle {
		bucketType = minio.BucketLookupPath
	}

	minioClient, err := minio.New(uri.Host, &minio.Options{
		Region:       s.config.Region,
		Creds:        credentials.NewStaticV4(s.config.AccessKey, s.config.SecretKey, ""),
		Secure:       uri.Scheme == "https",
		BucketLookup: bucketType,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}

func (s *S3Walker) Walk() error {
	if s.blockFlag {
		return nil
	}
	s.blockFlag = true

	s.Stats.Reset()
	s.startProcessing()
	buckets, err := s.client.ListBuckets(context.Background())

	if s.config.Bucket == "" {
		if err != nil {
			log.Errorf("Could not list buckets: %s", err)
		}

		for i := range buckets {
			if utils.MatchExclude(s.bucketPatterns, buckets[i].Name) {
				log.Infof("Bucket %s excluded !", buckets[i].Name)
				continue
			}
			s.walkBucket(context.Background(), buckets[i])
		}
	} else {
		s.walkBucket(context.Background(), minio.BucketInfo{Name: s.config.Bucket})
	}

	s.endProcessing()
	s.blockFlag = false

	return err
}

func (s *S3Walker) ValidateConfig(Config) error {
	return nil
}

func (s *S3Walker) walkBucket(context_bg context.Context, bucket minio.BucketInfo) {

	//ctxWithWait := context.WithValue(context_bg, "group", wait)
	s.findObjects(context_bg, bucket)
	//s.findIncompleteupload(ctxWithWait, bucket)

	log.Debug("Done listing objects for bucket ", bucket.Name)
}

func (s *S3Walker) findObjects(context_bg context.Context, bucket minio.BucketInfo) {
	//wait := context_bg.Value("group").(*sync.WaitGroup)
	ctx, cancel := context.WithCancel(context_bg)
	//defer wait.Done()
	defer cancel()

	objectCh := s.client.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			log.Warning("Object warning", object.Err.Error())
			continue
		}
		s.ProcessFile(bucket.Name,
			object.Key, object.Size,
			s.baseWalker.config.Depth,
			object.ContentType,
			map[string]string{"bucket": bucket.Name, "storageClass": object.StorageClass})
	}
}

//
//func (s *S3Walker) findIncompleteupload(context_bg context.Context, bucket minio.BucketInfo) {
//	wait := context_bg.Value("group").(*sync.WaitGroup)
//	ctx, cancel := context.WithCancel(context_bg)
//	defer wait.Done()
//	defer cancel()
//
//	objectCh := s.client.ListIncompleteUploads(ctx, bucket.Name, "", true)
//	for object := range objectCh {
//		if object.Err != nil {
//			fmt.Println(object.Err)
//			return
//		}
//		s.ProcessFile("", "/"+object.Key, object.Size, s.config.Depth, "")
//	}
//
//}
