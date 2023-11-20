locals {
  data = {
    moo = "cow",
    key = sensitive("super secret"),
  }
}

resource "null_resource" "data" {
  triggers = local.data
}

resource "terraform_data" "data" {
  input = local.data
}

resource "local_file" "this" {
  filename = "${path.module}/tmp/a.json"
  content  = jsonencode(local.data)
}
