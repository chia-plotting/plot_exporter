package plotinfo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

var (
	computingTableRegex      = regexp.MustCompile("^Computing table ([0-9]+)$")
	computingTableAltRegex   = regexp.MustCompile("^\\[P1\\] Table ([0-9]+) took .* sec(, found [0-9]+ matches)?$")
	backpropTableRegex       = regexp.MustCompile("^Backpropagating on table ([0-9]+)$")
	backpropTableAltRegex    = regexp.MustCompile("^\\[P2\\] Table ([0-9]+) rewrite took .* sec(, dropped .* entries \\(.*\\))?$")
	compressTableRegex       = regexp.MustCompile("^Compressing tables ([0-9]+) and ([0-9])+$")
	compressTableAltRegex    = regexp.MustCompile("^\\[P3-2\\] Table ([0-9]+) took .* sec, wrote [0-9]+ left entries, [0-9]+ final$")
	writeCheckpointTables    = regexp.MustCompile(".*Write Checkpoint tables.*$")
	writeCheckpointTablesAlt = regexp.MustCompile("^Total plot creation time was .* sec$")

	computingTableProgress = map[uint]uint{
		1: 1,
		2: 6,
		3: 12,
		4: 20,
		5: 28,
		6: 36,
		7: 42,
	}
	backpropTableProgress = map[uint]uint{
		2: 61,
		3: 58,
		4: 55,
		5: 51,
		6: 48,
		7: 43,
	}
	compressTableProgress = map[uint]uint{
		1: 66,
		2: 73,
		3: 79,
		4: 85,
		5: 92,
		6: 98,
		7: 99,
	}
)

/*
 * Returns the plot progress.
 *
 * Params:
 *   logReader - the reader of a log file of the chia plotter
 *
 * Returns:
 *   progress \in [0.100] - how close we are to completion
 *   completed - whether or not the plot is done
 *
 * See https://github.com/Chia-Network/chia-blockchain/wiki/Beginners-Guide#how-plots-are-created-and-7-steps-process for more details.
 */
func GetPlotProgress(logReader io.Reader) (progress uint, completed bool) {
	scanner := bufio.NewScanner(logReader)
	for scanner.Scan() && !completed {
		prog, done, err := getLineProgress(scanner.Text())
		if err != nil {
			log.Printf("GetPlotProgress(): %s", err)
		}

		if prog > progress {
			progress = prog
		}

		if done {
			completed = true
		}
	}
	return
}

func getLineProgress(line string) (uint, bool, error) {
	computingTableMatches := computingTableRegex.FindStringSubmatch(line)
	if len(computingTableMatches) == 0 {
		computingTableMatches = computingTableAltRegex.FindStringSubmatch(line)
	}

	if len(computingTableMatches) >= 2 {
		progressMeter, err := strconv.ParseUint(computingTableMatches[1], 10, 32)
		if err != nil {
			return 0, false, fmt.Errorf("getLineProgress(): %s", err)
		} else if progressMeter > 7 {
			return 0, false, fmt.Errorf("getLineProgress(): %d > 7", progressMeter)
		}
		return computingTableProgress[uint(progressMeter)], false, nil
	}

	backpropTableMatches := backpropTableRegex.FindStringSubmatch(line)
	if len(backpropTableMatches) == 0 {
		backpropTableMatches = backpropTableAltRegex.FindStringSubmatch(line)
	}

	if len(backpropTableMatches) >= 2 {
		progressMeter, err := strconv.ParseUint(backpropTableMatches[1], 10, 32)
		if err != nil {
			return 0, false, fmt.Errorf("getLineProgress(): %s", err)
		} else if progressMeter > 7 || progressMeter < 2 {
			return 0, false, fmt.Errorf("getLineProgress(): invalid table no: %d", progressMeter)
		}
		return backpropTableProgress[uint(progressMeter)], false, nil
	}

	compressTableMatches := compressTableRegex.FindStringSubmatch(line)
	if len(compressTableMatches) == 0 {
		compressTableMatches = compressTableAltRegex.FindStringSubmatch(line)
	}

	if len(compressTableMatches) >= 2 {
		progressMeter, err := strconv.ParseUint(compressTableMatches[1], 10, 32)
		if err != nil {
			return 0, false, fmt.Errorf("getLineProgress(): %s", err)
		} else if progressMeter > 6 {
			return 0, false, fmt.Errorf("getLineProgress(): invalid table no: %d", progressMeter)
		}
		return compressTableProgress[uint(progressMeter)], false, nil
	}

	if writeCheckpointTables.MatchString(line) || writeCheckpointTablesAlt.MatchString(line) {
		return 100, true, nil
	}

	return 0, false, nil
}
