local tempo = import 'github.com/grafana/tempo/operations/tempo-mixin/mixin.libsonnet';
local utils = import 'utils.libsonnet';

tempo {
  _config+:: {
  },
  grafanaDashboards+:: utils.withBrowserTimezone(super.grafanaDashboards),
}
