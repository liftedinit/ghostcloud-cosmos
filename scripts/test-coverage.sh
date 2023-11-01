#!/bin/bash
set -e -x

go test -v -coverprofile=coverage.txt -covermode=atomic -coverpkg=./x/... $(go list ./x/...)

# append "||true" to grep so if no match the return code stays 0
excludelist="$(find ./ -type f -name '*.go' | xargs grep -l 'DONTCOVER' || true)"
excludelist+=" $(find ./ -type f -name '*.pb.go')"
excludelist+=" $(find ./ -type f -name '*.pb.gw.go')"
excludelist+=" $(find ./ -type f -name '*_simulation.go')"
for filename in ${excludelist}; do
  filename=${filename#".//"}
  echo "Excluding ${filename} from coverage report..."
  filename=$(echo "$filename" | sed 's/\//\\\//g')
  sed -i.bak "/""$filename""/d" coverage.txt
done