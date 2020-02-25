#!/bin/bash

for dir in */cmd/*/
do
  cd "${dir}" || exit 1
  echo "cleaning ${dir}"
  go clean
  cd - > /dev/null  || exit 1
done