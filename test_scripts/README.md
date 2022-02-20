# Ecoji tests

This directory contains a test script that can be used to test other
implementations of Ecoji. To test another implemention you need to override what
commands are used to decode and encode. For example assume that `/home/ubuntu/bin/myEcoji`
is a new implementation executable.  To test this executable run the following
in this directory.

```bash
encode1_cmd="/home/ubuntu/bin/myEcoji -w 0" \ 
encode2_cmd="/home/ubuntu/bin/myEcoji -e -w 0" \
decode_cmd="/home/ubuntu/bin/myEcoji -d"  \
./ecoji_test.sh
```
