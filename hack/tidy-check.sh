#!/usr/bin/env bash
# Copyright 2019 Antrea Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
set +e
THIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
PROJECT_DIR="$THIS_DIR/.."
pushd "$THIS_DIR" >/dev/null || exit
MOD_FILE="$THIS_DIR/../go.mod"
SUM_FILE="$THIS_DIR/../go.sum"
TMP_DIR="$THIS_DIR/.tmp.tidy-check"
TMP_MOD_FILE="$TMP_DIR/go.mod"
TMP_SUM_FILE="$TMP_DIR/go.sum"
TARGET_GO_VER="go1.12.*"

general_help() {
  echo "Please run the following command to generate a new go.mod & go.sum:"
  if [[ "$(go version | awk '{print $3}')" == $TARGET_GO_VER ]]; then
    echo "  \$ make tidy"
  else
    echo "  \$ make docker-tidy"
  fi
}

precheck() {
  if [ ! -r "$MOD_FILE" ]; then
    echo "no go.mod found"
    general_help
    exit 1
  fi
  if [ ! -r "$SUM_FILE" ]; then
    echo "no go.sum found"
    general_help
    exit 1
  fi
  mkdir -p "$TMP_DIR"
}

tidy() {
  cp "$MOD_FILE" "$TMP_MOD_FILE"
  mv "$SUM_FILE" "$TMP_SUM_FILE"
  if [[ "$(go version | awk '{print $3}')" == $TARGET_GO_VER ]]; then
    go mod tidy >>/dev/null 2>&1
  else
    docker run -e "GOPROXY=$(go env GOPROXY)" -w /antrea -v "$PROJECT_DIR":/antrea golang:1.12 go mod tidy >>/dev/null 2>&1
  fi
}

clean() {
  mv "$TMP_MOD_FILE" "$MOD_FILE"
  mv "$TMP_SUM_FILE" "$SUM_FILE"
  rm -fr "$TMP_DIR"
}

failed() {
  echo "'go mod tidy' failed, there are errors in dependencies rules"
  general_help
  clean
  exit 1
}

check() {
  MOD_DIFF=$(diff "$MOD_FILE" "$TMP_MOD_FILE")
  SUM_DIFF=$(diff "$SUM_FILE" "$TMP_SUM_FILE")
  if [ -n "$MOD_DIFF" ] || [ -n "$SUM_DIFF" ]; then
    echo "dependencies are not tidy"
    general_help
    clean
    exit 1
  fi
  clean
}

precheck
if tidy; then
  check
else
  failed
fi

popd >/dev/null || exit
