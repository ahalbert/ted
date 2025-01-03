#!/bin/zsh

set -o errexit
set -o nounset
set -o pipefail

for testfile in tests/**/*.fsa; do
  testname=$(basename $testfile | sed 's/.fsa$//')
  echo "running test $testfile..."
  infile=$(echo $testfile | sed 's/.fsa$/.in/')
  outfile=$(echo $testfile | sed 's/.fsa$/.out/')
  flags=$(cat "$testfile:A:h/flags")
  ./bin/ted -f "$testfile" $(echo $flags) "$infile" > ./bin/output
  if ! diff ./bin/output "$outfile" > /dev/null; then
    echo "ERROR: test $testname failed!"
  fi
done
