# Copyright (c) 2021 Terminus, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script is used to generate the pb.go file for the extension package.
# The generated pb.go file is used to register the extension package to the
# global registry.

set -e
set -o pipefail

cd "$(dirname "$0")"
## at tools dir
cd ../../../../../tools
protoc --go_out=protoc/include --go_opt=paths=source_relative -I=protoc/include  protoc/include/custom/extension/extension.proto
echo done
