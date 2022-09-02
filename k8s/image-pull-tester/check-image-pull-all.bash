#!/bin/bash

check-image-pull() {
  local ctx="$1"
  echo "================"
  echo "check $ctx image pull"
  kubectl --context "$ctx" apply -f image-pull-tester.yaml
  sleep 10
  echo "----------------"
  kubectl --context "$ctx" -n image-pull-tester get pod -owide
  echo "----------------"
  kubectl --context "$ctx" delete ns image-pull-tester
  echo "================"
}

check-image-pull fc-com-aks
check-image-pull fc-dev-aks
check-image-pull fc-personal-aks
check-image-pull fc-stg-aks
