local loki = import 'loki-mixin/mixin.libsonnet';

loki {
  _config+:: {
    promtail+: {
      enabled: false,
    },
    thanos+: {
      enabled: false,
    },
    operational+: {
      memcached: false,
      consul: false,
      bigTable: false,
      dynamo: false,
      gcs: false,
      s3: true,
      azureBlob: false,
      boltDB: false,
    },
    ssd+: {
      enabled: true,
    },
  },
}
