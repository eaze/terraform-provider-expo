terraform {
  required_providers {
    expo = {
      source = "eaze/expo"
    }
  }
}

provider "expo" {
  token        = "..."
  account_name = "..."
}
