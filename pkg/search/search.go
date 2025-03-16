package search

func (cse *CentaureissiSearchEngine) Index(msg CentaureissiSearchDocument) error {
	err := cse.index.Index(msg.Hash, msg)
	if err != nil {
		return err
	}
	return nil
}
