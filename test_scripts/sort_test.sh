#!/bin/bash

#export LC_ALL=en_US.utf8
export LC_ALL=C

rm /tmp/sort-test.ecoji

# TODO need to generate a more exhaustive sort test

for i in $(seq 0 255); do 
  printf "%02x" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
# begin test of 4 special padding chars
  printf "%02x00 0000" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0000 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0000 01" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0000 ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0001" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0001 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0001 ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0002" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0002 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0002 ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0003" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0003 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0003 ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 0004" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
# end test of 4 special padding chars
  printf "%02x00 00ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 01" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x00 ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
# the first two bits of 2nd byte influence 1st output emoji, the follwing four go through all 2 bit combos
  printf "%02x01" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x40" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02x80" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xc0" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff 01" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ff00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ff01" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ffff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ffff 00" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ffff 01" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
  printf "%02xff ffff ff" $i | xxd -r -p | ecoji >> /tmp/sort-test.ecoji
done

sort /tmp/sort-test.ecoji > /tmp/sort-sorted.ecoji
diff /tmp/sort-test.ecoji /tmp/sort-sorted.ecoji
