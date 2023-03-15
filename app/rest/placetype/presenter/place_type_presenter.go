package placetype_presenter

type PlaceTypePresenterCTX interface {
	JSON(code int, i interface{}) error
	Bind(i interface{}) error
	Param(name string) string
}
