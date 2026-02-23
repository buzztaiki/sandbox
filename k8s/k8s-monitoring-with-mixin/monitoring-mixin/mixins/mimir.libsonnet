local utils = import './lib/utils.libsonnet';
local mimir = import 'github.com/grafana/mimir/operations/mimir-mixin/mixin.libsonnet';

mimir {
  _config+:: {
  },
  grafanaDashboards+:: utils.withBrowserTimezone(super.grafanaDashboards),
}
