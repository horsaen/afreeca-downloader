package tools

type ConcurrentRow struct {
	Index  int
	Values []string
}

func SnapshotConcurrentRow(index int, values []string) ConcurrentRow {
	snapshot := make([]string, len(values))
	copy(snapshot, values)

	return ConcurrentRow{
		Index:  index,
		Values: snapshot,
	}
}
