package dogs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhargreaves/dog-walking/test/acceptance/common"
)

type Dog struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createDog(t *testing.T, name string) Dog {
	resp := common.PostJson(t, "/dogs", Dog{Name: name})
	defer resp.Body.Close()
	common.RequireStatus(t, resp, http.StatusCreated)

	var dog Dog
	common.DecodeJSON(t, resp, &dog)

	assert.Equal(t, name, dog.Name, "Expected dog name to be returned")
	assert.NotEmpty(t, dog.ID, "Expected dog ID to be returned")

	return dog
}
