#!/bin/bash
go build -o terraform-provider-openwhisk .

rm out.zip 2>/dev/null
cd example
zip -r ../out.zip *
cd ..

terraform init
terraform apply