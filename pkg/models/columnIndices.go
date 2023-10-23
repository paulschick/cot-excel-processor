package models

type ColumnNameIndices struct {
	BaseOrderedColumnNames []string
	FullOrderedColumnNames []string
	ColumnNamesIndex       map[string]int
}

func NewColumnNameIndices() *ColumnNameIndices {
	return &ColumnNameIndices{
		BaseOrderedColumnNames: []string{
			MarketExchangeName,
			ReportDate,
			OpenInterest,
			NonCommLong,
			NonCommShort,
			CommLong,
			CommShort,
			NonReptLong,
			NonReptShort,
		},
		FullOrderedColumnNames: []string{
			MarketExchangeName,
			ReportDate,
			OpenInterest,
			NonCommLong,
			NonCommShort,
			CommLong,
			CommShort,
			NonReptLong,
			NonReptShort,
			NonCommNet,
			CommNet,
			NonReptNet,
		},
		ColumnNamesIndex: make(map[string]int),
	}
}

func (c *ColumnNameIndices) InitializeColumnIndices(headerRow []string) {
	for _, colName := range c.BaseOrderedColumnNames {
		c.ColumnNamesIndex[colName] = -1
	}
	for colIdx := 0; colIdx < len(headerRow); colIdx++ {
		col := headerRow[colIdx]
		if _, ok := c.ColumnNamesIndex[col]; ok {
			c.ColumnNamesIndex[col] = colIdx
		}
	}
}
