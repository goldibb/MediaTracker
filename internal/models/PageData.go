package models

type PageData struct {
	PageTitle  string
	PageMode   string
	Groups     []GroupedItems
	SortStyles []string
}
