package api

/*
Basic API Response wrapper
All API queries will satisfy this struct at least
*/
type Response struct {
	Status string `xml:"status,attr"`
}

/*
Entry represents any PANOS <entry> object and the methods it requires.
*/
type Entry interface {
	Print()
}
