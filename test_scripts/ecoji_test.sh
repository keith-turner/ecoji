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

echo "INFO checking Ecoji V2 concatenated data"
if ! diff <(echo -n "👖📸🧈🌭👩☕💲🥇🪚☕" | $decode_cmd) <(echo -n "abcdefxyz") &> /dev/null; then
	echo "ERROR failed to decode Ecoji V2 concatenated data"
fi

echo "INFO checking Ecoji V1 concatenated data"
if ! diff <(echo -n "👖📸🎈☕🎥🤠📠🏍🐲👡🕟☕" | $decode_cmd) <(echo "abc6789XY") &> /dev/null; then
	echo "ERROR failed to decode Ecoji V1 concatenated data"
fi

