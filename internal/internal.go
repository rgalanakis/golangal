package internal

import "github.com/onsi/gomega"

func CoerceToMatcher(i interface{}) gomega.OmegaMatcher {
	if m, ok := i.(gomega.OmegaMatcher); ok {
		return m
	}
	return gomega.Equal(i)
}
