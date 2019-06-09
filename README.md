# spanner-operator

Spanner operation client for Cloud Spanner Administrator.  
This repo is based on [sample-controller](https://github.com/kubernetes/sample-controller).  
This is an experimental project of my own practice creating custom controller.

# Plans

- Implement custom metrics server

# References

## About custom controller

- [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/)
- [sample-controller](https://github.com/kubernetes/sample-controller)
- [code-generator](https://github.com/kubernetes/code-generator)
- [client-go](https://github.com/kubernetes/client-go)
- [oracle/mysql-operator](https://github.com/oracle/mysql-operator)

## About autoscaling

- [Horizontal Pod Autoscaler Walkthrough](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)
- [podautoscaler/horizontal.go](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/podautoscaler/horizontal.go)
- [Autoscaling Deployments with External Metrics](https://cloud.google.com/kubernetes-engine/docs/tutorials/external-metrics-autoscaling)
- [Autoscaling Deployments with Custom Metrics](https://cloud.google.com/kubernetes-engine/docs/tutorials/custom-metrics-autoscaling)

## About custom metrics

- [Configure the Aggregation Layer](https://kubernetes.io/docs/tasks/access-kubernetes-api/configure-aggregation-layer/)
- [Kubernetes Autoscaling with Custom Metrics](https://www.infracloud.io/kubernetes-autoscaling-custom-metrics/)
- [GCP Metrics List](https://cloud.google.com/monitoring/api/metrics_gcp)
- [Reading Metric Data](https://cloud.google.com/monitoring/custom-metrics/reading-metrics?hl=en#monitoring_read_timeseries_fields-go)
- [k8s-stackdriver/custom-metrics-stackdriver-adapter](https://github.com/GoogleCloudPlatform/k8s-stackdriver/tree/master/custom-metrics-stackdriver-adapter)
- [DirectXMan12/k8s-prometheus-adapter](https://github.com/DirectXMan12/k8s-prometheus-adapter)
