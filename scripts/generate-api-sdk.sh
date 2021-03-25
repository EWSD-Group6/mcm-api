#!/usr/bin/env bash
docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i /local/docs/swagger.json \
    -g typescript-angular \
    -o /local/api-module