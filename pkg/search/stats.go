package search

func (cse *CentaureissiSearchEngine) StatsDocuments() uint64 {
	docCount, _ := cse.index.DocCount()
	return docCount
}
