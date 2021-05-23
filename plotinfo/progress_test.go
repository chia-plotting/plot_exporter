package plotinfo

import (
	"strings"
	"testing"
)

func TestGetPlotProgressEmpty(t *testing.T) {
	sr := strings.NewReader("")
	p, c := GetPlotProgress(sr)
	if p != 0 {
		t.Errorf("progress = %d, expected progress = 0", p)
	} else if c {
		t.Errorf("completed = %v, expected completed = false", c)
	}
}

func TestGetPlotProgressRoutine(t *testing.T) {
	sr := strings.NewReader(`Computing table 7
Some Noisy Line
Compressing tables 6 and 7`)
	p, c := GetPlotProgress(sr)
	if p != 98 {
		t.Errorf("progress = %d, expected progress = 98", p)
	} else if c {
		t.Errorf("completed = %v, expected completed = false", c)
	}
}

func TestGetLineProgressEmpty(t *testing.T) {
	s := ""
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 0 {
		t.Errorf("getLineProgress(\"%s\") returned a non-zero progress", s)
	}
}

func TestGetLineProgressRoutineComputingTable1(t *testing.T) {
	s := "Computing table 1"
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 1 {
		t.Errorf("getLineProgress(\"%s\") returned progress %d, expected %d", s, p, 1)
	}
}

func TestGetLineProgressRoutineComputingTable7(t *testing.T) {
	s := "Computing table 7"
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 42 {
		t.Errorf("getLineProgress(\"%s\") returned progress %d, expected %d", s, p, 42)
	}
}

func TestGetLineProgressComputingTableInvalid(t *testing.T) {
	s := "Computing table 8"
	_, _, e := getLineProgress(s)
	if e == nil {
		t.Errorf("getLineProgress(\"%s\") did not return an error", s)
	}
}

func TestGetLineProgressRoutineBackpropagatingTable7(t *testing.T) {
	s := "Backpropagating on table 7"
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 43 {
		t.Errorf("getLineProgress(\"%s\") returned progress %d, expected %d", s, p, 43)
	}
}

func TestGetLineProgressRoutineBackpropagatingTable2(t *testing.T) {
	s := "Backpropagating on table 2"
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 61 {
		t.Errorf("getLineProgress(\"%s\") returned progress %d, expected %d", s, p, 61)
	}
}

func TestGetLineProgressInvalidBackpropagatingTable8(t *testing.T) {
	s := "Backpropagating on table 8"
	_, _, e := getLineProgress(s)
	if e == nil {
		t.Errorf("getLineProgress(\"%s\") did not return an error", s)
	}
}

func TestGetLineProgressInvalidBackpropagatingTable1(t *testing.T) {
	s := "Backpropagating on table 1"
	_, _, e := getLineProgress(s)
	if e == nil {
		t.Errorf("getLineProgress(\"%s\") did not return an error", s)
	}
}

func TestGetLineProgressRoutineCompressingTables1(t *testing.T) {
	s := "Compressing tables 1 and 2"
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 66 {
		t.Errorf("getLineProgress(\"%s\") returned progress %d, expected %d", s, p, 66)
	}
}

func TestGetLineProgressRoutineCompressingTables6(t *testing.T) {
	s := "Compressing tables 6 and 7"
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if d {
		t.Errorf("getLineProgress(\"%s\") returned done expectedly", s)
	} else if p != 98 {
		t.Errorf("getLineProgress(\"%s\") returned progress %d, expected %d", s, p, 66)
	}
}

func TestGetLineProgressInvalidCompressingTables7(t *testing.T) {
	s := "Compressing tables 7 and 8"
	_, _, e := getLineProgress(s)
	if e == nil {
		t.Errorf("getLineProgress(\"%s\") did not return an error", s)
	}
}

func TestGetLineProgressWriteCheckpointTables(t *testing.T) {
	s := `Starting phase 4/4: Write Checkpoint tables into "/mount/ssd0/plot-k32-2021-05-23-09-43-abf3373e4cf81ac349ef2302255d713a8357fc34191fabc38af013fe0941e858.plot.2.tmp" ... Sun May 23 14:26:38 2021`
	p, d, e := getLineProgress(s)
	if e != nil {
		t.Errorf("getLineProgress(\"%s\") returned an unexpected error", s)
	} else if !d {
		t.Errorf("getLineProgress(\"%s\") did not return done", s)
	} else if p != 100 {
		t.Errorf("getLineProgress(\"%s\") did not return progress = 100", s)
	}
}
