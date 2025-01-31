#!/usr/bin/bash
operating_systems=(linux windows)
version=0.4

for os in ${operating_systems[@]}
do
    env GOOS=${os} GOARCH=amd64 go build -o bin/go-secret-app-${version}-${os}-amd64.bin main.go
    
    if [ "windows" == "${os}" ]; then 
        mv bin/go-secret-app-${version}-${os}-amd64.bin bin/go-secret-app-${version}-${os}-amd64.exe

        go-base64 encode bin/go-secret-app-${version}-${os}-amd64.exe bin/go-secret-app-${version}-${os}-amd64.exe-base64.txt
        go-base64 decode bin/go-secret-app-${version}-${os}-amd64.exe-base64.txt bin/go-secret-app-${version}-${os}-amd64.exe-base64.txt_decoded.bin
        md5sum bin/go-secret-app-${version}-${os}-amd64.exe
        md5sum bin/go-secret-app-${version}-${os}-amd64.exe-base64.txt_decoded.bin

    fi 

done

ls -alh bin/ 

