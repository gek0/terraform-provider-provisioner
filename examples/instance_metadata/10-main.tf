provider "provisioner" {  
  # initialize provider with:
  # api_endpoint or export PROVISIONER_API_ENDPOINT (optional since default is set)
  # api_key or export PROVISIONER_API_KEY
}

resource "provisioner_instance_metadata" "instance_metadata" {
  ad_domain     = var.ad_domain
  instance_name = var.instance_name
}

output "instance_metadata_id" {
  value = provisioner_instance_metadata.instance_metadata.id
}
