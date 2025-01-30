#!/usr/bin/bash
operating_systems=(linux windows)
version=0.2

for os in ${operating_systems[@]}
do
    env GOOS=${os} GOARCH=amd64 go build -o bin/go-secret-app-${version}-${os}-amd64.bin main.go
    
    if [ "windows" == "${os}" ]; then 
        mv bin/go-secret-app-${version}-${os}-amd64.bin bin/go-secret-app-${version}-${os}-amd64.exe
    fi 

done

ls -alh bin/ 

