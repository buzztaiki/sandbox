local kubernetes = import 'github.com/kubernetes-monitoring/kubernetes-mixin/mixin.libsonnet';
local utils = import 'utils.libsonnet';

kubernetes {
  _config+:: {
    showMultiCluster: true,
  },
  grafanaDashboards: utils.withBrowserTimezone({
    [x.key]: x.value
    for x in std.objectKeysValues(super.grafanaDashboards)
    if !(std.startsWith(x.key, 'k8s-resources-windows') || std.startsWith(x.key, 'k8s-windows'))
  }),
}
