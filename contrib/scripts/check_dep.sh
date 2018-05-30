#!/bin/bash

expected="context
github.com/lehajam/dgo/protos/api
github.com/lehajam/dgo/y
github.com/gogo/protobuf/proto
github.com/pkg/errors
google.golang.org/grpc/codes
google.golang.org/grpc/status
math/rand
sync"

got=$(go list -f '{{ join .Imports "\n" }}' github.com/lehajam/dgo)

if [ "$got" != "$expected" ]; then
  echo "Imports for Go client didn't match."
  echo -e "\nExpected, Got\n"
  diff -y <(echo "$expected") <(echo "$got")
  exit 1
fi
