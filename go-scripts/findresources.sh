#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 <resource_kind> <desired_status>"
  exit 1
fi

resource_kind=$1
desired_status=$2

for ns in $(kubectl get namespaces -o=jsonpath='{range .items[*].metadata.name} {.} {"\n"} {end}'); do
  echo "Checking resources of kind $resource_kind in namespace $ns"
  resources=$(kubectl get $resource_kind -n $ns -o=jsonpath='{range .items[*].metadata.name} {.} {"\n"} {end}')

  for resource in $resources; do
    current_status=$(kubectl get $resource_kind $resource -n $ns -o=jsonpath='{.status.phase}')

    if [ "$current_status" != "$desired_status" ]; then
      echo "$resource in namespace $ns is not in the desired status ($desired_status)"
    fi
  done
done
