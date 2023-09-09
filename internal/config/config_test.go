package config

import "testing"

func Test_Config_Validate(t *testing.T) {
	emptyConfig := NewConfiguration()

	err := emptyConfig.Validate()
	if err == nil {
		t.Errorf("Error validating config: %v", err)
	}
}
