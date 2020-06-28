resource "openwhisk_function" "openwhisk_test" {
    name = "test"
    zip_path = "out.zip"
    environment = {
        testvar = "First test function"
    }
}

resource "openwhisk_function" "openwhisk_other_test" {
    name = "test2"
    zip_path = "out.zip"
    environment = {
        testvar = "Second test function"
    }
}