#!/usr/bin/env bash

USERNAME=${1:-codespace}
FAILED_COUNT=0

if [ -z $HOME ]; then
  HOME="/root"
fi

function echoStderr() {
  echo "$@" 1>&2
}

function check() {
  LABEL=$1
  shift
  echo -e "\n🧪 Testing: ${LABEL}"
  if "$@"; then
    echo "✅ Passed: ${LABEL}"
    return 0
  else
    echoStderr "❌ Failed: ${LABEL}"
    FAILED_COUNT=$((FAILED_COUNT + 1))
    return 1
  fi
}

function reportResults() {
  if [ ${FAILED_COUNT} -eq 0 ]; then
    hundred="💯"
    printf "\n%s All tests passed!\n" "${hundred}"
    exit 0
  else
    boom="💥"
    printf "\n%s Failed tests: %s\n" "${boom}" "${FAILED_COUNT}"
    exit 1
  fi
}

function non_root_user() {
  id ${USERNAME}
}
