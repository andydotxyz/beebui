#!/bin/sh

DIR=`dirname "$0"`
FILE=bundled.go
BIN=`go env GOPATH`/bin

cd $DIR

$BIN/fyne bundle -package beebui -name monitor monitor.png > $FILE

$BIN/fyne bundle -package beebui -name font -append kongtext.ttf >> $FILE
