package model

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultErrorInternal(t *testing.T) {
	fn := model.FlowNetwork{}
	assert.True(t, model.IsFNCreator(&fn))
	assert.True(t, model.IsFNCreator(fn))
}
