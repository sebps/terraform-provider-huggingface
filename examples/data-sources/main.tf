terraform {
  required_providers {
    huggingface = {
      source   = "hashicorp.com/edu/huggingface"
      version  = "~> 1.0"
    }
  }
}

provider "huggingface" {
  hf_token = "<YOUR_HF_TOKEN>"
}

data "huggingface_endpoints" "example" {
  namespace = "<YOUR_NAMESPACE>"
}

output "endpoints" {
  value = data.huggingface_endpoints.example
}