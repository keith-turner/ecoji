#!/bin/bash

: ${encode1_cmd:="ecoji -w 0"}
: ${encode2_cmd:="ecoji -e -w 0"}
: ${decode_cmd:="ecoji -d"}

for plain_file in data/*.plain; do
	bare_name=$(basename $plain_file .plain)

	echo "INFO checking $bare_name"

	if ! diff <(cat data/$bare_name.ev1) <(cat $plain_file | $encode1_cmd) &> /dev/null; then
		echo "ERROR : Encoding $plain_file using Ecojiv1 did not produce data/$bare_name.ev1"
	fi

	if ! diff <(cat data/$bare_name.ev2) <(cat $plain_file | $encode2_cmd) &> /dev/null; then
		echo "ERROR : Encoding $plain_file using Ecojiv2 did not produce data/$bare_name.ev2"
	fi

	if !  diff <(cat $plain_file) <(cat data/$bare_name.ev1 | $decode_cmd) &> /dev/null; then
		echo "ERROR : Decoding data/$bare_name.ev1 did not produce $plain_file"
	fi

	if ! diff <(cat $plain_file) <(cat data/$bare_name.ev2 | $decode_cmd) &> /dev/null; then
		echo "ERROR : Decoding data/$bare_name.ev2 did not produce $plain_file"
	fi
done

#Run test that only decode data
for plain_file in data/*.plaind; do
	bare_name=$(basename $plain_file .plaind)

	echo "INFO checking $bare_name"

	if !  diff <(cat $plain_file) <(cat data/$bare_name.enc | $decode_cmd) &> /dev/null; then
		echo "ERROR : Decoding data/$bare_name.enc did not produce $plain_file"
	fi
done


#Run test that only decode data
for plain_file in data/*.garbage; do
	bare_name=$(basename $plain_file .garbage)

	echo "INFO checking $bare_name"

	if  cat data/$bare_name.garbage | $decode_cmd &> /dev/null; then
		echo "ERROR : Decoding data/$bare_name.garbage did not produce an error code"
	fi
done

# TODO this script does not test wrapping in any way

