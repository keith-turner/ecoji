package ecoji

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestMakrdown(t *testing.T) {

	f, err := os.Create("docs/emojis.md")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	fmt.Fprintln(writer, "# Emojis used by Ecoji")

	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, " * [Emojis that differ between Ecoji V1 and V2](#emojis-that-differ-between-ecoji-v1-and-v2)")
	fmt.Fprintln(writer, " * [Emojis that are same in Ecoji V1 and V2](#emojis-that-are-same-in-ecoji-v1-and-v2)")
	fmt.Fprintln(writer, " * [Emojis used for padding](#emojis-used-for-padding)")
	fmt.Fprintln(writer, " * [Candidate not used in Ecoji V2](#candidate-not-used-in-ecoji-v2)")
	fmt.Fprintln(writer, " * [Emojis in Ecoji V1 not in candidates](#emojis-in-ecoji-v1-not-in-candidates)")
	fmt.Fprintln(writer, "")

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "<!-- This file is automatically generated -->")
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "## Emojis that differ between Ecoji V1 and V2")
	fmt.Fprintln(writer)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "The candidate column indicates if a V1 emoji was a [candidate](candidates.md) for V2.  If a V1 emojis was a candidate and was not used in V2 then someone decided against using it in V2.  If a V1 emoji was not a candidate, then it could not be used in V2 because it did not meet the selection criteria.")
	fmt.Fprintln(writer)

	candidateRunes := getRunes("docs/candidates.txt")

	candidatesMap := make(map[rune]bool)
	for _, r := range candidateRunes {
		candidatesMap[r] = true
	}

	fmt.Fprintln(writer, "Ordinal | V1 codepoint | V1 emoji | candidate | V2 codepoint | V2 emoji")
	fmt.Fprintln(writer, "-|-|-|-|-")
	for i, r := range emojisV1 {
		if emojisV2[i] != r {
			candidate := "Y"
			if _, present := candidatesMap[r]; !present {
				candidate = "N"
			}
			fmt.Fprintf(writer, "%d | %U | %s | %s | %U | %s\n", i, r, string(r), candidate, emojisV2[i], string(emojisV2[i]))
		}
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "## Emojis that are same in Ecoji V1 and V2")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Ordinal | codepoint | emoji ")
	fmt.Fprintln(writer, "-|-|-")
	for i, r := range emojisV1 {
		if emojisV2[i] == r {
			fmt.Fprintf(writer, "%d | %U | %s \n", i, r, string(r))
		}
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "## Emojis used for padding")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Type | codepoint V1 | emoji V1 | codepoint V1 | emoji V2 ")
	fmt.Fprintln(writer, "-|-|-|-|-")

	fmt.Fprintf(writer, "FILL | %U | %s | %U | %s\n", padding, string(padding), padding, string(padding))
	for i, r := range paddingLastV1 {
		fmt.Fprintf(writer, "PAD_%d | %U | %s | %U | %s\n", i, r, string(r), paddingLastV2[i], string(paddingLastV2[i]))
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "## Candidate not used in Ecoji V2")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "The following are [candidates](candidates.md) that were not used. This information is provided for reference and is not needed to implement Ecoji")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "codepoint | emoji ")
	fmt.Fprintln(writer, "-|-")
	for _, r := range candidateRunes {
		if _, present := revEmojis[r]; !present {
			fmt.Fprintf(writer, "%U | %s \n", r, string(r))
		}
	}
}

func getRunes(fileName string) []rune {

	var runes []rune

	file, err := os.Open(fileName)
	handle(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		i, err := strconv.ParseInt(scanner.Text(), 16, 32)
		handle(err)
		runes = append(runes, rune(i))
	}

	handle(scanner.Err())

	return runes
}

func handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
