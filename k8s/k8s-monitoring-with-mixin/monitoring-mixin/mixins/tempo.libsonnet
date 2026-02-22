local tempo = import 'tempo-mixin/mixin.libsonnet';
local utils = import 'utils.libsonnet';

tempo {
  _config+:: {
  },
  grafanaDashboards+:: utils.withBrowserTimezone(super.grafanaDashboards),
}
