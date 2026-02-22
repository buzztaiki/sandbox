local mimir = import 'mimir-mixin/mixin.libsonnet';

mimir {
  _config+:: {
  },
  grafanaDashboards+:: {
    [filename]+: { timezone: 'browser' }
    for filename in std.objectFields(super.grafanaDashboards)
  },
}
