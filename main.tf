locals {
  fns = {
    "test"  = "./example/build/out.zip"
    "test2" = "./example/build/out.zip"
  }
}

resource "openwhisk_function" "fns" {
  for_each = local.fns
  name     = each.key
  source   = each.value
  environment = {
    FAASTERMETRICS_DEPLOYMENT_ID = "456deadbeef"
//    TIMESTAMP                    = timestamp()
  }
}
