#!/bin/bash
rm test.exe
rm blockChain.db
rm blockChain.db.lock
go build -o "test.exe"