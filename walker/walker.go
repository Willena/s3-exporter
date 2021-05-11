package walker

import log "github.com/sirupsen/logrus"

type Config struct {
	*BaseWalkerConfig
	*S3WalkerConfig
	*FsWalkerConfig
}

type Walker interface {
	Init(config Config, labels map[string]string, labelValue []string) error
	Walk() error
	ValidateConfig(config Config) error
}

func FromConfig(config Config, walkerType string) (Walker, error) {
	var walker Walker

	switch walkerType {
	case "s3":
		walker = &S3Walker{}

	case "fs":
		walker = &FsWalker{}
	default:
		log.Fatalln("Please select one walker implementation")
		return nil, nil
	}

	err := walker.Init(config, config.CustomLabels, nil)

	return walker, err

}
