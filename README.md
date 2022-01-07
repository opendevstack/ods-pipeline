# Sample Helm chart

You may use this as a starting point for a sample application deployed with the `ods-deploy-helm` task.

## Usage

Download this as a ZIP file, either via "Code > Download ZIP" or via:

```
curl -LO https://github.com/opendevstack/ods-pipeline/archive/refs/heads/sample-helm-chart.zip
unzip sample-helm-chart.zip
mv ods-pipeline-sample-helm-chart chart
rm sample-helm-chart.zip
```

Afterwards, adjust the sample Helm chart to your needs. In any case, you must edit `Chart.yaml` and change the values of `name` and `description` to fit your component, e.g. for a repository named `foo-bar`, the chart's `name` should be `bar`.
