#!/bin/bash

scripts/git.sh --add_commit
scripts/git.sh --add_push

cd ../

mkdir -p target
mkdir -p tmp

go mod tidy
go generate