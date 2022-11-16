#!/bin/bash

if [[ $# -ne 1 ]]; then
  echo "usage: $1 <context>" 1>&2
  exit 1
fi

context="$1"

# shellcheck disable=SC2016
query='
.items[]
  | . as $res
  | [$res.spec.template, $res.spec.jobTemplate.spec.template][] | select(. != null)
  | [.spec.initContainers, .spec.containers][] | select(. != null) | flatten[]
  | [$res.metadata.namespace, $res.kind, $res.metadata.name, .image]
  | @tsv
'
kubectl --context "$context" get deploy,statefulset,daemonset,cronjob -ojson -A | jq "$query" -r -M
