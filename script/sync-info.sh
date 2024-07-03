#!/bin/bash
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


# 设置项目根目录
project_root="."

# 遍历项目目录下的所有文件夹
for dir in $(find "$project_root" -type d); do
    if [[ "$dir" == *node_modules* ]] || [[ "$dir" == *.git* ]] || [[ "$dir" == *.vscode* ]]; then
        continue
    fi

    if [ -f "$dir/info.yaml" ]; then
        version=$(awk '/version:/{print $2}' "$dir/info.yaml")

        if [ -f "$dir/package.json" ]; then
            sed -i '' -E 's/"version": "[0-9]+\.[0-9]+\.[0-9]+"/"version": "'$version'"/' "$dir/package.json"
        fi
    fi
done
