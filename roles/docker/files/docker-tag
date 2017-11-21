#!/bin/bash
curl -s -S "https://registry.hub.docker.com/v2/repositories/$@/tags/" | jq '."results"[]["name"]' |sort
