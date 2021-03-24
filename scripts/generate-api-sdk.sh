#!/usr/bin/env bash
java -jar ./scripts/swagger-codegen-cli.jar generate -i docs/swagger.json -l typescript-angular -o api-module