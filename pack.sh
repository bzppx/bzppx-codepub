#!/bin/bash
VER="0.8"
RELEASE="release-${VER}"
rm -rf ${RELEASE}
mkdir ${RELEASE}

# windows amd64
echo 'Start pack windows amd64...'
cd install
GOOS=windows GOARCH=amd64 go build ./
cd ..
GOOS=windows GOARCH=amd64 go build ./
tar -czvf "${RELEASE}/bzppx-codepub-windows-amd64.tar.gz" bzppx-codepub.exe conf/ static/ views/ install/install.exe LICENSE README.md
rm -rf bzppx-codepub.exe

echo 'Start pack windows X386...'
cd install
GOOS=windows GOARCH=386 go build ./
cd ..
GOOS=windows GOARCH=386 go build ./

tar -czvf "${RELEASE}/bzppx-codepub-windows-386.tar.gz" bzppx-codepub.exe conf/ static/ views/ install/install.exe LICENSE README.md
rm -rf bzppx-codepub.exe

echo 'Start pack linux amd64'
cd install
GOOS=linux GOARCH=amd64 go build ./
cd ..
GOOS=linux GOARCH=amd64 go build ./
tar -czvf "${RELEASE}/bzppx-codepub-linux-amd64.tar.gz" bzppx-codepub conf/ static/ views/ install/install LICENSE README.md
rm -rf bzppx-codepub

echo 'Start pack linux 386'
cd install
GOOS=linux GOARCH=386 go build ./
cd ..
GOOS=linux GOARCH=386 go build ./

tar -czvf "${RELEASE}/bzppx-codepub-linux-386.tar.gz" bzppx-codepub conf/ static/ views/ install/install LICENSE README.md
rm -rf bzppx-codepub

echo 'END'
