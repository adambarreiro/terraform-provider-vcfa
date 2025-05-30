---
page_title: "VMware Cloud Foundation Automation: vcfa_edge_cluster_qos"
subcategory: ""
description: |-
  Provides a VMware Cloud Foundation Automation Edge Cluster QoS data source.
---

# vcfa_edge_cluster_qos

Provides a resource to read Quality of Service (QoS) settings from an Edge Cluster in VMware Cloud Foundation Automation.

_Used by: **Provider**_

## Example Usage

```hcl
data "vcfa_region" "demo" {
  name = "region-one"
}

data "vcfa_edge_cluster" "demo" {
  name             = "my-edge-cluster"
  region_id        = data.vcfa_region.demo.id
  sync_before_read = true
}

data "vcfa_edge_cluster_qos" "demo" {
  edge_cluster_id = data.vcfa_edge_cluster.demo.id
}
```

## Argument Reference

The following arguments are supported:

- `edge_cluster_id` - (Required) An ID of Edge Cluster. Can be looked up using
  [vcfa_edge_cluster](/providers/vmware/vcfa/latest/docs/data-sources/edge_cluster) data source

## Attribute Reference

All the arguments and attributes defined in
[`vcfa_edge_cluster_qos`](/providers/vmware/vcfa/latest/docs/resources/edge_cluster_qos) resource are available.
