#!/usr/bin/env bash
usage() {
    if [ "$err" != "" ]; then
        echo $err
        echo ""
    fi

    echo "$0 [path to ci.yml]"
    echo ""
    echo "Attempt to retrieve the slack notification channel from the CI configuration."
    echo "This is a helper for the Makefile's \`staging\` and \`prod\` targets."
    echo ""
    exit 1
}

ci_file=${1:-.release/ci.hcl}
if [ "$#" -gt 1 ]; then
    err="ERROR: too many parameters."
    usage
fi

if [ ! -f ${ci_file} ]; then
    err="ERROR: File not found - ${ci_file}."
    usage
fi

channel="$(cat "${ci_file}" | awk '/notification_channel/{print $3}' | tr -d '"' )"
echo "${channel}"