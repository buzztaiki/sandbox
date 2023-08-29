local g = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local azureMonitorDS = 'grafana-azure-monitor-datasource';

local appInsightPanel = function(panelName, resourceGroup, resourceName) (
  local target = function(metrics) (
    // local azureMonitor = g.query.azureMonitor.azureMonitor;
    // g.query.azureMonitor.withSubscription('$subscription')
    // + azureMonitor.withAggregation('Total')
    // + azureMonitor.withDimensionFilters([])
    // + azureMonitor.withMetricName(metrics)
    // + azureMonitor.withMetricNamespace('microsoft.insights/components')
    // + azureMonitor.withResourceGroup(resourceGroup)
    // + azureMonitor.withResourceName(resourceName)
    // + azureMonitor.withTimeGrain('auto')
    // + azureMonitor.withTop('')

    {
      azureMonitor: {
        aggregation: 'Total',
        dimensionFilters: [],
        metricName: metrics,
        metricNamespace: 'microsoft.insights/components',
        resourceGroup: resourceGroup,
        resourceName: resourceName,
        timeGrain: 'auto',
        top: '',
      },
      subscription: '$subscription',
    }
  );

  g.panel.timeSeries.new(panelName)
  + g.panel.timeSeries.datasource.withType(azureMonitorDS)
  + g.panel.timeSeries.datasource.withUid('$datasource')
  + g.panel.timeSeries.withTargets([
    target('requests/count'),
    target('requests/failed'),
    target('exceptions/count'),
  ])
);


local variables = [
  g.dashboard.variable.datasource.new('datasource', azureMonitorDS),
  g.dashboard.variable.query.new('subscription', 'Subscriptions()')
  + g.dashboard.variable.query.withDatasource(type=azureMonitorDS, uid='$datasource'),
];


local panels = (
  local appInsighs = import 'app-insights.json';
  g.util.grid.makeGrid([
    g.panel.row.new(rg)
    + g.panel.row.withPanels([appInsightPanel(name, rg, name) for name in appInsighs[rg]])
    for rg in std.objectFields(appInsighs)
  ], panelWidth=6, panelHeight=6)
);


local main = (
  g.dashboard.new('Azure All App Insights')
  + g.dashboard.withTimezone('browser')
  + g.dashboard.withUid('7ab5509e-697d-44f6-a851-9d10fdbfbfd1')
  + g.dashboard.withVariables(variables)
  + g.dashboard.graphTooltip.withSharedCrosshair()
  + g.dashboard.withPanels(panels)
);

main
