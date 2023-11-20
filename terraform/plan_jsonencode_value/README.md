# jsonencode された値の差分を plan の時に見る

`null_resource` を使ったハック。https://github.com/hashicorp/terraform/issues/28947#issuecomment-1323451595 で紹介されていたもの。`terraform_data` だと、https://github.com/hashicorp/terraform/issues/32789 にあるように `output` で全部漏れちゃってよろしくない。

## first plan

```
❯❯ tf plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # local_file.this will be created
  + resource "local_file" "this" {
      + content              = (sensitive value)
      + content_base64sha256 = (known after apply)
      + content_base64sha512 = (known after apply)
      + content_md5          = (known after apply)
      + content_sha1         = (known after apply)
      + content_sha256       = (known after apply)
      + content_sha512       = (known after apply)
      + directory_permission = "0777"
      + file_permission      = "0777"
      + filename             = "./tmp/a.json"
      + id                   = (known after apply)
    }

  # null_resource.data will be created
  + resource "null_resource" "data" {
      + id       = (known after apply)
      + triggers = {
          + "key" = (sensitive value)
          + "moo" = "cow"
        }
    }

  # terraform_data.data will be created
  + resource "terraform_data" "data" {
      + id     = (known after apply)
      + input  = {
          + key = (sensitive value)
          + moo = "cow"
        }
      + output = (known after apply)
    }

Plan: 3 to add, 0 to change, 0 to destroy.
```

## edit
```
diff --git a/terraform/plan_jsonencode_value/main.tf b/terraform/plan_jsonencode_value/main.tf
index 61ee193..512598b 100644
--- a/terraform/plan_jsonencode_value/main.tf
+++ b/terraform/plan_jsonencode_value/main.tf
@@ -1,6 +1,6 @@
 locals {
   data = {
-    moo = "cow",
+    moo = "bull",
     key = sensitive("super secret"),
   }
 }
```

## second plan
```
❯❯ tf plan
terraform_data.data: Refreshing state... [id=7575bd0d-524b-2c6d-bea9-3c78bf4ab3b0]
null_resource.data: Refreshing state... [id=5521072920525248069]
local_file.this: Refreshing state... [id=78c006d482032702ca03ee68cc314316ad74d975]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # local_file.this must be replaced
-/+ resource "local_file" "this" {
      ~ content              = (sensitive value) # forces replacement
      ~ content_base64sha256 = "l8zw1M47cDHaXyiBUj2ICpFqdRq6tBZBj9qseEprBaQ=" -> (known after apply)
      ~ content_base64sha512 = "W+pAp/byqfqlFwpaiBpBuLsub4GDApt/8ghaRqttOhbNnBNi3iXxqvceL97v1Ii74RXgjqwfOFRuCotSHhlsmQ==" -> (known after apply)
      ~ content_md5          = "1dcc6206c43027787c260a2fe582b884" -> (known after apply)
      ~ content_sha1         = "78c006d482032702ca03ee68cc314316ad74d975" -> (known after apply)
      ~ content_sha256       = "97ccf0d4ce3b7031da5f2881523d880a916a751abab416418fdaac784a6b05a4" -> (known after apply)
      ~ content_sha512       = "5bea40a7f6f2a9faa5170a5a881a41b8bb2e6f8183029b7ff2085a46ab6d3a16cd9c1362de25f1aaf71e2fdeefd488bbe115e08eac1f38546e0a8b521e196c99" -> (known after apply)
      ~ id                   = "78c006d482032702ca03ee68cc314316ad74d975" -> (known after apply)
        # (3 unchanged attributes hidden)
    }

  # null_resource.data must be replaced
-/+ resource "null_resource" "data" {
      ~ id       = "5521072920525248069" -> (known after apply)
      ~ triggers = { # forces replacement
          ~ "moo" = "cow" -> "bull"
            # (1 unchanged element hidden)
        }
    }

  # terraform_data.data will be updated in-place
  ~ resource "terraform_data" "data" {
        id     = "7575bd0d-524b-2c6d-bea9-3c78bf4ab3b0"
      ~ input  = {
          ~ moo = "cow" -> "bull"
            # (1 unchanged attribute hidden)
        }
      ~ output = {
          - key = "super secret"
          - moo = "cow"
        } -> (known after apply)
    }

Plan: 2 to add, 1 to change, 2 to destroy.
```
