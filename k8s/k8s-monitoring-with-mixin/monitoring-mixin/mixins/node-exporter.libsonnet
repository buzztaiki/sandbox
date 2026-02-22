local node = import 'node-mixin/mixin.libsonnet';

node {
  _config+:: {
    nodeExporterSelector: 'job="node-exporter"',
  },
  grafanaDashboards: {
    [x.key]: x.value
    for x in std.objectKeysValues(super.grafanaDashboards)
    if !(std.startsWith(x.key, 'nodes-aix') || std.startsWith(x.key, 'nodes-darwin'))
  },
}
