package models

import "strconv"

const (
	MarketExchangeName = "Market_and_Exchange_Names"
	ReportDate         = "Report_Date_as_MM_DD_YYYY"
	OpenInterest       = "Open_Interest_All"
	NonCommLong        = "NonComm_Positions_Long_All"
	NonCommShort       = "NonComm_Positions_Short_All"
	CommLong           = "Comm_Positions_Long_All"
	CommShort          = "Comm_Positions_Short_All"
	NonReptLong        = "NonRept_Positions_Long_All"
	NonReptShort       = "NonRept_Positions_Short_All"
	NonCommNet         = "NonComm_Positions_Net_All"
	CommNet            = "Comm_Positions_Net_All"
	NonReptNet         = "NonRept_Positions_Net_All"
)

type Row struct {
	MarketExchangeName string
	ReportDate         string
	OpenInterest       string
	NonCommLong        int
	NonCommShort       int
	NonCommNet         int
	CommLong           int
	CommShort          int
	CommNet            int
	NonReptLong        int
	NonReptShort       int
	NonReptNet         int
}

func NewRow(row []string, colNameIndex map[string]int) *Row {
	r := &Row{
		MarketExchangeName: row[colNameIndex[MarketExchangeName]],
		ReportDate:         row[colNameIndex[ReportDate]],
		OpenInterest:       row[colNameIndex[OpenInterest]],
		NonCommLong:        parseInt(row[colNameIndex[NonCommLong]]),
		NonCommShort:       parseInt(row[colNameIndex[NonCommShort]]),
		CommLong:           parseInt(row[colNameIndex[CommLong]]),
		CommShort:          parseInt(row[colNameIndex[CommShort]]),
		NonReptLong:        parseInt(row[colNameIndex[NonReptLong]]),
		NonReptShort:       parseInt(row[colNameIndex[NonReptShort]]),
	}
	r.parseNetValues()
	return r
}

func (r *Row) GetCsvRow() []string {
	return []string{
		r.MarketExchangeName,
		r.ReportDate,
		r.OpenInterest,
		strconv.Itoa(r.NonCommLong),
		strconv.Itoa(r.NonCommShort),
		strconv.Itoa(r.CommLong),
		strconv.Itoa(r.CommShort),
		strconv.Itoa(r.NonReptLong),
		strconv.Itoa(r.NonReptShort),
		strconv.Itoa(r.NonCommNet),
		strconv.Itoa(r.CommNet),
		strconv.Itoa(r.NonReptNet),
	}
}

func (r *Row) parseNetValues() {
	r.NonCommNet = r.NonCommLong - r.NonCommShort
	r.CommNet = r.CommLong - r.CommShort
	r.NonReptNet = r.NonReptLong - r.NonReptShort
}

func parseInt(value string) int {
	if value, err := strconv.Atoi(value); err == nil {
		return value
	}
	return 0
}
