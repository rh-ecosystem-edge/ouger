#!/bin/bash

update_k8s_api() {
  local latest_tag=$(curl -s https://api.github.com/repos/kubernetes/api/git/refs/tags | jq -r 'map(.ref | select(test("refs/tags/v[0-9]+\\.[0-9]+\\.[0-9]+$"))) | map(sub("refs/tags/";"")) | last')

  if [[ -z "${latest_tag}" ]]; then
    echo "Error: could not fetch the latest release version for Kubernetes."
    exit 1
  fi

  echo "Latest stable version of k8s.io/api: ${latest_tag}"

  go get k8s.io/api@"${latest_tag}"
  go get k8s.io/apimachinery@"${latest_tag}"
  go get k8s.io/kube-aggregator@"${latest_tag}"
  go get k8s.io/kubectl@"${latest_tag}"
}

update_openshift_api() {
  local latest_tag=$(curl -s https://api.github.com/repos/openshift/api/branches | jq -r 'map(select(.name | test("^release-4\\.[0-9]+$"))) | map(.name) | last')

  if [[ -z "${latest_tag}" ]]; then
    echo "Error: could not fetch the latest release version for OpenShift."
    exit 1
  fi

  echo "Latest stable version of openshift/api: ${latest_tag}"

  go get "github.com/openshift/api@${latest_tag}"
}

main() {
 update_k8s_api

 update_openshift_api
}

main "$@"
