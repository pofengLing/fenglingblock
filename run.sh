#!/bin/bash
rm test.exe
rm blockChain.db
rm blockChain.db.lock
rm wallet.dat
go build -o "test.exe"