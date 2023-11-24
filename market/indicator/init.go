package indicator

func LoadIndicators() {
	// load all the indicators
	// Note: this is a hack to make sure all the indicators are loaded as they utilize the init
	// function and go only calls it if something from this package is called.
	// Maybe someday this will actually do something
}
