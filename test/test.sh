#!/usr/bin/env bash

SCRIPT_FOLDER_NAME=$(dirname $0)
cd ${SCRIPT_FOLDER_NAME} || exit

# shellcheck source=/dev/null
source test-utils.sh

check "hadolint" hadolint --version
check "shfmt" shfmt --version
check "shellcheck" shellcheck --version
check "editorconfig" ec --version
check "trivy" trivy --version

reportResults
