package reflection_test

import (
	"testing"

	"github.com/monadicstack/abide/internal/reflection"
	"github.com/stretchr/testify/suite"
)

func TestReflectionSuite(t *testing.T) {
	suite.Run(t, new(ReflectionSuite))
}

type ReflectionSuite struct {
	suite.Suite
}

func (suite *ReflectionSuite) TestToBindingValue() {
	r := suite.Require()

	type organization struct {
		ID   int
		Name string
	}

	type group struct {
		ID     int
		Name   string
		Org    organization `json:"Organization"`
		OrgPtr *organization
	}

	type user struct {
		ID    int
		Name  string `json:"alias"`
		Group group
	}

	// empty := user{}
	dude := user{
		ID:   123,
		Name: "Dude",
		Group: group{
			ID:   456,
			Name: "Bowling League",
			Org: organization{
				ID:   789,
				Name: "Lebowski-Fest",
			},
			OrgPtr: &organization{
				ID:   42,
				Name: "His Dudeness",
			},
		},
	}

	testIntValue := func(u user, path string, expected int) {
		var intValue int
		r.True(reflection.ToBindingValue(u, path, &intValue))
		r.Equal(expected, intValue)
	}
	testStringValue := func(u user, path string, expected string) {
		var stringValue string
		r.True(reflection.ToBindingValue(u, path, &stringValue))
		r.Equal(expected, stringValue)
	}

	var intValue int
	var stringValue string

	// Garbage data tests
	r.False(reflection.ToBindingValue(dude, "Turds", &intValue))
	r.False(reflection.ToBindingValue(dude, "", &intValue))
	r.False(reflection.ToBindingValue(stringValue, "ID", &intValue))
	r.False(reflection.ToBindingValue(dude, "ID", &stringValue))
	r.Panics(func() {
		reflection.ToBindingValue(dude, "ID", nil)
	})

	// Can properly fetch primitive fields at the root level
	testIntValue(dude, "ID", 123)
	testStringValue(dude, "alias", "Dude")

	// Can go recursively deep for values
	testIntValue(dude, "Group.ID", 456)
	testStringValue(dude, "Group.Name", "Bowling League")
	testIntValue(dude, "Group.Organization.ID", 789)
	testIntValue(dude, "Group.OrgPtr.ID", 42)
	testStringValue(dude, "Group.OrgPtr.Name", "His Dudeness")

	// Can grab complex data structures as binding values.
	var groupValue group
	r.True(reflection.ToBindingValue(dude, "Group", &groupValue))
	r.Equal(456, groupValue.ID)
	r.Equal("Bowling League", groupValue.Name)

	var orgValue organization
	r.True(reflection.ToBindingValue(dude, "Group.Organization", &orgValue))
	r.Equal(789, orgValue.ID)
	r.Equal("Lebowski-Fest", orgValue.Name)

	var orgPtrValue *organization
	r.True(reflection.ToBindingValue(dude, "Group.OrgPtr", &orgPtrValue))
	r.Equal(42, orgPtrValue.ID)
	r.Equal("His Dudeness", orgPtrValue.Name)

	// When a field is remapped using the `json` tag, you need to use that in the binding path. The actual
	// field name should NOT work!
	r.False(reflection.ToBindingValue(dude, "Name", &stringValue))
	r.False(reflection.ToBindingValue(dude, "Group.Org", &intValue))
	r.False(reflection.ToBindingValue(dude, "Group.Org.ID", &intValue))
}
