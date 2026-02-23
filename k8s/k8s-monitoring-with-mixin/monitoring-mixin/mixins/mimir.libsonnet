local mimir = import 'github.com/grafana/mimir/operations/mimir-mixin/mixin.libsonnet';
local utils = import './lib/utils.libsonnet';

mimir {
  _config+:: {
  },
  grafanaDashboards+:: utils.withBrowserTimezone(super.grafanaDashboards),
}
