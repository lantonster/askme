#!/bin/bash

generate_mocks() {
    local source_dir=$1
    local destination_dir=$2

    files=$(find "$source_dir" -type f -name "*.go")

    for file in $files
    do
        base=$(basename "$file")
        name="${base%.*}"
        # 获取文件所在目录相对于 source_dir 的路径作为子目录部分
        relative_dir=$(dirname "$file" | sed "s|$source_dir/||") 
        if [ $relative_dir = $source_dir ]; then
            destination="$destination_dir/${name}_mock.go"
        else 
            destination="$destination_dir/$relative_dir/${name}_mock.go"
        fi
        # 获取文件所在目录名作为包名
        package_name=$(dirname "$file")
        package_name=$(basename "$package_name")
        mockgen -source="$file" -destination="$destination" -package="$package_name"
    done
}

if [ $# -ne 2 ]; then
    echo "Usage: $0 <source_directory> <destination_directory>"
    exit 1
fi

generate_mocks "$1" "$2"