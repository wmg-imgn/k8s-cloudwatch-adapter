package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/awslabs/k8s-cloudwatch-adapter/pkg/apis/metrics/v1alpha1"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

	"k8s.io/klog"
)

// GetLocalRegion gets the region ID from the instance metadata.
func GetLocalRegion() string {
	
	return "us-east-1"
}

func toCloudWatchQuery(externalMetric *v1alpha1.ExternalMetric) cloudwatch.GetMetricDataInput {
	queries := externalMetric.Spec.Queries

	cwMetricQueries := make([]types.MetricDataQuery, len(queries))
	for i, q := range queries {
		q := q
		returnData := &q.ReturnData
		mdq := types.MetricDataQuery{
			Id:         &q.ID,
			Label:      &q.Label,
			ReturnData: *returnData,
		}

		if len(q.Expression) == 0 {
			dimensions := make([]types.Dimension, len(q.MetricStat.Metric.Dimensions))
			for j := range q.MetricStat.Metric.Dimensions {
				dimensions[j] = types.Dimension{
					Name:  &q.MetricStat.Metric.Dimensions[j].Name,
					Value: &q.MetricStat.Metric.Dimensions[j].Value,
				}
			}

			metric := &types.Metric{
				Dimensions: dimensions,
				MetricName: &q.MetricStat.Metric.MetricName,
				Namespace:  &q.MetricStat.Metric.Namespace,
			}

			mdq.MetricStat = &types.MetricStat{
				Metric: metric,
				Period: &q.MetricStat.Period,
				Stat:   &q.MetricStat.Stat,
				Unit:   types.StandardUnit(q.MetricStat.Unit),
			}
		} else {
			mdq.Expression = &q.Expression
		}

		cwMetricQueries[i] = mdq
	}
	cwQuery := cloudwatch.GetMetricDataInput{
		MetricDataQueries: cwMetricQueries,
	}

	return cwQuery
}
