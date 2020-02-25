#!/bin/bash

for dir in */cmd/*/
do
  cd "${dir}" || exit 1
  echo "building ${dir}"
  go clean && go build -o main
  cd - > /dev/null  || exit 1
done