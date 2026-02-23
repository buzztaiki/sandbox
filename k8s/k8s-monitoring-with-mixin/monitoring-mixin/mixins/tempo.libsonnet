local utils = import './lib/utils.libsonnet';
local tempo = import 'github.com/grafana/tempo/operations/tempo-mixin/mixin.libsonnet';

tempo {
  _config+:: {
  },
  grafanaDashboards+:: utils.withBrowserTimezone(super.grafanaDashboards),
}
