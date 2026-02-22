local node = import 'node-mixin/mixin.libsonnet';

node {
  _config+:: {
    nodeExporterSelector: 'job="node-exporter"',
  },
}
