local utils = import './lib/utils.libsonnet';
local kubernetes = import 'github.com/kubernetes-monitoring/kubernetes-mixin/mixin.libsonnet';

kubernetes {
  _config+:: {
    showMultiCluster: true,
    common_join_labels: [
      'label_app_kubernetes_io_name',
      'label_app_kubernetes_io_instance',
      'label_app_kubernetes_io_component',
    ],
  },
  grafanaDashboards: utils.withBrowserTimezone({
    [x.key]: x.value
    for x in std.objectKeysValues(super.grafanaDashboards)
    if !(std.startsWith(x.key, 'k8s-resources-windows') || std.startsWith(x.key, 'k8s-windows'))
  }),
}
