package locationtype_presenter

type LocationTypePresenterCTX interface {
	JSON(code int, i interface{}) error
	Bind(i interface{}) error
	Param(name string) string
}
