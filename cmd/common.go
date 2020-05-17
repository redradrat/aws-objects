package cmd

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

func GetSession(cmd *cobra.Command) (client.ConfigProvider, error) {
	reg, err := cmd.Flags().GetString(RegionFlag)
	if err != nil {
		return nil, err
	}
	conf := &awssdk.Config{}
	if reg != "" {
		conf.Region = awssdk.String(reg)
	}

	cp, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	return cp, nil
}
