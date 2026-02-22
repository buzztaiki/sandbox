local kubernetes = import 'kubernetes-mixin/mixin.libsonnet';

kubernetes {
  _config+:: {
    grafanaK8s+: {
      grafanaTimezone: 'browser',
    },
    showMultiCluster: true,
  },
  grafanaDashboards: {
    [x.key]: x.value
    for x in std.objectKeysValues(super.grafanaDashboards)
    if !(std.startsWith(x.key, 'k8s-resources-windows') || std.startsWith(x.key, 'k8s-windows'))
  },
}
