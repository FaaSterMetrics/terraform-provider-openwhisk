#!/bin/bash

set -euo pipefail

go build -o terraform-provider-openwhisk .

rm -f example/build/out.zip
cd example
npm run build
cd build
zip -r out.zip *
cd ../../

terraform init
terraform apply
