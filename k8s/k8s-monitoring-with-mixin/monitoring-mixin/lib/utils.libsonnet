{
  // withBrowserTimezone sets timezone to 'browser' for all dashboards in the given object.
  withBrowserTimezone(dashboards):
    { [name]+: { timezone: 'browser' } for name in std.objectFields(dashboards) },
}
