package modules

import "fmt"

// newWideDemoTable returns many columns and a few rows (horizontal overflow demo).
func newWideDemoTable() (colKeys, colTitles []string, data [][]string) {
	const n = 22
	colKeys = make([]string, n)
	colTitles = make([]string, n)
	for i := 0; i < n; i++ {
		colKeys[i] = fmt.Sprintf("k%02d", i+1)
		colTitles[i] = fmt.Sprintf("col_%02d", i+1)
	}
	const numRows = 6
	data = make([][]string, numRows)
	for r := 0; r < numRows; r++ {
		row := make([]string, n)
		for c := 0; c < n; c++ {
			row[c] = fmt.Sprintf("v%d:%d", r+1, c+1)
		}
		data[r] = row
	}
	return colKeys, colTitles, data
}

// newTallDemoTable returns a few columns and many rows (vertical pagination demo).
func newTallDemoTable() (colKeys, colTitles []string, data [][]string) {
	colKeys = []string{"id", "status", "qty", "note", "day"}
	colTitles = []string{"id", "status", "qty", "note", "day"}
	statuses := []string{"queued", "run", "done", "fail", "hold"}
	const numRows = 92
	data = make([][]string, 0, numRows)
	for i := 1; i <= numRows; i++ {
		data = append(data, []string{
			fmt.Sprintf("%d", i),
			statuses[i%len(statuses)],
			fmt.Sprintf("%d", (i*17)%999),
			fmt.Sprintf("row %d payload", i),
			fmt.Sprintf("2026-%02d-%02d", 1+(i%4), 1+(i%27)),
		})
	}
	return colKeys, colTitles, data
}
