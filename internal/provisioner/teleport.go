// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package provisioner

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattermost/mattermost-cloud/internal/tools/aws"
	"github.com/mattermost/mattermost-cloud/internal/tools/kops"
	"github.com/mattermost/mattermost-cloud/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type teleport struct {
	awsClient      aws.AWS
	environment    string
	provisioner    *KopsProvisioner
	kops           *kops.Cmd
	cluster        *model.Cluster
	logger         log.FieldLogger
	desiredVersion string
	actualVersion  string
}

func newTeleportHandle(cluster *model.Cluster, desiredVersion string, provisioner *KopsProvisioner, awsClient aws.AWS, kops *kops.Cmd, logger log.FieldLogger) (*teleport, error) {
	if logger == nil {
		return nil, errors.New("cannot instantiate Teleport handle with nil logger")
	}

	if provisioner == nil {
		return nil, errors.New("cannot create a connection to Teleport if the provisioner provided is nil")
	}

	if kops == nil {
		return nil, errors.New("cannot create a connection to Teleport if the Kops command provided is nil")
	}

	environment, err := getEnvironment(awsClient)
	if err != nil {
		return nil, err
	}

	if environment == "" {
		return nil, errors.New("cannot create a connection to Teleport if the environment is empty")
	}

	return &teleport{
		awsClient:      awsClient,
		environment:    environment,
		provisioner:    provisioner,
		kops:           kops,
		cluster:        cluster,
		logger:         logger.WithField("cluster-utility", model.TeleportCanonicalName),
		desiredVersion: desiredVersion,
	}, nil

}

func (n *teleport) updateVersion(h *helmDeployment) error {
	actualVersion, err := h.Version()
	if err != nil {
		return err
	}

	n.actualVersion = actualVersion
	return nil
}

func (n *teleport) CreateOrUpgrade() error {
	h := n.NewHelmDeployment()
	err := h.Update()
	if err != nil {
		return err
	}

	err = n.updateVersion(h)
	return err
}

func (n *teleport) DesiredVersion() string {
	return n.desiredVersion
}

func (n *teleport) ActualVersion() string {
	return strings.TrimPrefix(n.actualVersion, "teleport-")
}

func (n *teleport) Destroy() error {
	teleportClusterName := fmt.Sprintf("cloud-%s-%s", n.environment, n.cluster.ID)
	err := n.awsClient.S3EnsureBucketDeleted(teleportClusterName, n.logger)
	if err != nil {
		return errors.Wrap(err, "unable to delete Teleport bucket")
	}

	err = n.awsClient.DynamoDBEnsureTableDeleted(teleportClusterName, n.logger)
	if err != nil {
		return errors.Wrap(err, "unable to delete Teleport dynamodb table")
	}

	err = n.awsClient.DynamoDBEnsureTableDeleted(fmt.Sprintf("%s-events", teleportClusterName), n.logger)
	if err != nil {
		return errors.Wrap(err, "unable to delete Teleport dynamodb events table")
	}
	return nil
}

func (n *teleport) NewHelmDeployment() *helmDeployment {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = aws.DefaultAWSRegion
	}
	teleportClusterName := fmt.Sprintf("cloud-%s-%s", n.environment, n.cluster.ID)
	return &helmDeployment{
		chartDeploymentName: "teleport",
		chartName:           "chartmuseum/teleport",
		namespace:           "teleport",
		setArgument:         fmt.Sprintf("config.auth_service.cluster_name=%[1]s,config.teleport.storage.region=%[2]s,config.teleport.storage.table_name=%[1]s,config.teleport.storage.audit_events_uri=dynamodb://%[1]s-events,config.teleport.storage.audit_sessions_uri=s3://%[1]s/records?region=%[2]s", teleportClusterName, awsRegion),
		valuesPath:          "helm-charts/teleport_values.yaml",
		kopsProvisioner:     n.provisioner,
		kops:                n.kops,
		logger:              n.logger,
		desiredVersion:      n.desiredVersion,
	}
}

func (n *teleport) Name() string {
	return model.TeleportCanonicalName
}

func getEnvironment(awsClient aws.AWS) (string, error) {
	accountAliases, err := awsClient.GetAccountAliases()
	if err != nil {
		return "", err
	}
	if len(accountAliases.AccountAliases) < 1 {
		return "", errors.New("Account Alias not defined")
	}
	for _, alias := range accountAliases.AccountAliases {
		if strings.HasPrefix(*alias, "mattermost-cloud") && len(strings.Split(*alias, "-")) == 3 {
			return strings.Split(*alias, "-")[2], nil
		}
	}
	return "", errors.New("Account environment was not obtained")
}
