local node = import 'github.com/prometheus/node_exporter/docs/node-mixin/mixin.libsonnet';
local utils = import './lib/utils.libsonnet';

node {
  _config+:: {
    nodeExporterSelector: 'job="node-exporter"',
    showMultiCluster: true,
  },
  grafanaDashboards: utils.withBrowserTimezone({
    [x.key]: x.value
    for x in std.objectKeysValues(super.grafanaDashboards)
    if !(std.startsWith(x.key, 'nodes-aix') || std.startsWith(x.key, 'nodes-darwin'))
  }),
}
