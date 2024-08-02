package sanitizer

func MaskSecrets(m map[string]interface{}) func(i interface{}) error {
	return func(i interface{}) error {
		for k := range m {
			m[k] = "***"
		}

		return nil
	}
}
