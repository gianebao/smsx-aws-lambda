package model_test

import (
	"testing"

	"github.com/gianebao/smsx-aws-lambda/app/model"
	"github.com/stretchr/testify/assert"
)

func TestStatus_String(t *testing.T) {
	assert.Equal(t, `{"code":0,"text":"OK","network":100}`, model.Status{Code: 0,
		Text:    "OK",
		Network: 100}.String())
}
