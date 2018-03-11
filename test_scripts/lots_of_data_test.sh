#!/bin/bash

head -c 1000000000 /dev/urandom > /tmp/test.bin

sha1sum /tmp/test.bin
ecoji /tmp/test.bin | ecoji -d | sha1sum

rm /tmp/test.bin
