local tempo = import 'tempo-mixin/mixin.libsonnet';

tempo {
  _config+:: {
  },
  grafanaDashboards+:: {
    [filename]+: { timezone: 'browser' }
    for filename in std.objectFields(super.grafanaDashboards)
  },
}
