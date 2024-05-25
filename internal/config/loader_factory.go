package config

func CreateLoader(storage configStorage) *Loader {
	return &Loader{
		storage: storage,
		parsers: map[string]Parser{
			"yaml": NewYAMLParser(),
			"yml":  NewYAMLParser(),
		},
	}
}
