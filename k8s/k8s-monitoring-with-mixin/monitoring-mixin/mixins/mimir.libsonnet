local mimir = import 'mimir-mixin/mixin.libsonnet';
local utils = import 'utils.libsonnet';

mimir {
  _config+:: {
  },
  grafanaDashboards+:: utils.withBrowserTimezone(super.grafanaDashboards),
}
