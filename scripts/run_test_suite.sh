#!/usr/bin/env bash

set -euo pipefail

build/vc build tests/stdlib -o test-scripts
build/vc build tests/compiler -o test-scripts

failed=0

for script in test-scripts/*.vco; do
    [[ -f "$script" ]] || continue

    if ! build/vc run "$script"; then
        failed=1
    fi
done

source ./tests/bootstrap_test.sh
source ./tests/fixed_point_test.sh

exit "$failed"
