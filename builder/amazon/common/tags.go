package common

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type TagMap map[string]string
type EC2Tags []*ec2.Tag

func (t EC2Tags) Report(ui packer.Ui) {
	for _, tag := range t {
		ui.Message(fmt.Sprintf("Adding tag: \"%s\": \"%s\"",
			aws.StringValue(tag.Key), aws.StringValue(tag.Value)))
	}
}

func (t TagMap) IsSet() bool {
	return len(t) > 0
}

func (t TagMap) EC2Tags(ctx interpolate.Context, region, sourceAMIID string) (EC2Tags, error) {
	var ec2Tags []*ec2.Tag
	ctx.Data = &BuildInfoTemplate{
		SourceAMI:   sourceAMIID,
		BuildRegion: region,
	}
	for key, value := range t {
		interpolatedKey, err := interpolate.Render(key, &ctx)
		if err != nil {
			return nil, fmt.Errorf("Error processing tag: %s:%s - %s", key, value, err)
		}
		interpolatedValue, err := interpolate.Render(value, &ctx)
		if err != nil {
			return nil, fmt.Errorf("Error processing tag: %s:%s - %s", key, value, err)
		}
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   aws.String(interpolatedKey),
			Value: aws.String(interpolatedValue),
		})
	}
	return ec2Tags, nil
}
