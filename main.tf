resource "openwhisk_function" "openwhisk_test" {
    name = "test"
    zip_path = "./example/build/out.zip"
    environment = {
        FAASTERMETRICS_DEPLOYMENT_ID = "123deadbeef"
    }
}

resource "openwhisk_function" "openwhisk_other_test" {
    name = "test2"
    zip_path = "./example/build/out.zip"
    environment = {
        FAASTERMETRICS_DEPLOYMENT_ID = "456deadbeef"
    }
}
