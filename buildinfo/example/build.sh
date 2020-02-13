#!/usr/bin/env bash

set -x

# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=`git log --pretty=oneline -n 1`

# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}

# 检查源码在git commit 基础上，是否有本地修改，且未提交的内容
GitStatus=`git status -s`

# 获取当前时间
BuildTime=`date +'%Y.%m.%d.%H%M%S'`

# 获取 Go 的版本
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'github.com/rfyiamcool/golib/buildinfo.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/rfyiamcool/golib/buildinfo.GitStatus=${GitStatus}' \
    -X 'github.com/rfyiamcool/golib/buildinfo.BuildTime=${BuildTime}' \
    -X 'github.com/rfyiamcool/golib/buildinfo.BuildGoVersion=${BuildGoVersion}' \
"

ROOT_DIR=`pwd`

cd ${ROOT_DIR} && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/runner

echo 'build done.'
