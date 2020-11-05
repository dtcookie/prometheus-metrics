# Prometheus Metric Importer for Dynatrace

Dynatrace can automatically [monitor Prometheus Exporters](https://www.dynatrace.com/support/help/technology-support/cloud-platforms/kubernetes/monitoring/monitor-prometheus-metrics/) that are deployed in Kubernetes.

This project aims to take care of Prometheus Exporters located in any other infrastructure.
Instead of scraping various exporters actively, Prometheus itself is getting queried for any metrics it has digested so far.

Metrics are however just getting imported into Dynatrace, in case a correlation between the instance of the Prometheus metric and a known Host in Dynatrace can get found.
By default an import of the selected metrics is being performed every 50 seconds.

## Usage

You may specify configuration settings directly via command line arguments

```
Usage of ingest:
  -baseurl string
        The base url of your Managed Dynatrace Server
  -config string
        a JSON file containing settings
  -environment string
        the environment id of your Dynatrace Tenant
  -prometheus string
        Host and Port of your Prometheus Server
  -token string
        an API Token to access the REST API of your Dynatrace Tenant
```

or rely in full on a configuration file (default name `settings.json`).
A configuration file is necessary, because in here you can configure which metrics you would like to import and which ones not.

This sample for a configuration JSON file preconfigures all the necessary settings already. Placing a file name `settings.json` with these contents alongside with the `ingest` executable makes sure that you don't have to specify any additional arguments.

```
{
    "dynatrace": {
        "environment": "abc12345",
        "token": "##############"
    },
    "prometheus": "prometheus.local-network.org:9090",
    "metrics": {
        "includes": [
            "node_memory_*"
        ],
        "excludes": []
    }
}
```

The `includes` section needs to contain at least one entry. The utility refuses to import bluntly all available metrics.

## Visualization

All imported metrics are available via [Metric Explorer](https://www.dynatrace.com/support/help/how-to-use-dynatrace/dashboards-and-charts/explorer/) in Dynatrace for charting.

## Limitations

The importer currently just supports importing metrics that produce `vector` measurements. Other metrics are getting ignored.