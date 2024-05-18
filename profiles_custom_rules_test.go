package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListProfileCustomRules(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"rules": [
				  {
					"PK": "@RU",
					"order": 4,
					"group": 0,
					"action": {
					  "do": 0,
					  "status": 1
					}
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := ListProfileCustomRulesParams{
		ProfileID: "profileID",
		FolderID:  "folderID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/rules/%s", params.ProfileID, params.FolderID), handler)
	actual, err := client.ListProfileCustomRules(context.Background(), params)

	want := []Rule{
		{
			PK:    "@RU",
			Order: 4,
			Group: 0,
			Action: Action{
				Do:     0,
				Status: IntBool{true},
			},
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListProfileCustomRules(context.Background(), ListProfileCustomRulesParams{
		ProfileID: "",
		FolderID:  "folderID",
	})
	assert.Error(t, err, "Profile Rule Folders should not have been listed")
	_, err = client.ListProfileCustomRules(context.Background(), ListProfileCustomRulesParams{
		ProfileID: "profileID",
		FolderID:  "",
	})
	assert.Error(t, err, "Profile Rule Folders should not have been listed")
}

func TestCreateProfileCustomRule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"rules": [
				  {
					"do": 1,
					"status": 1,
					"order": 1
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := CreateProfileCustomRuleParams{
		ProfileID: "profileID",
		Do:        Bypass,
		Status:    IntBool{true},
		Hostnames: []string{"hostname1"},
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/rules", params.ProfileID), handler)
	actual, err := client.CreateProfileCustomRule(context.Background(), params)

	order := int64(1)
	want := []CustomRule{
		{
			Do:     Bypass,
			Status: IntBool{true},
			Order:  &order,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.CreateProfileRuleFolder(context.Background(), CreateProfileRuleFolderParams{
		ProfileID: "",
	})
	assert.Error(t, err, "Profile Custom Rule should not have been created")
}

func TestUpdateProfileCustomRule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"rules": [
				  {
					"do": 0,
					"status": 1,
					"order": 1,
					"group": 0
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := UpdateProfileCustomRuleParams{
		ProfileID: "profileID",
		Do:        Block,
		Status:    IntBool{true},
		Hostnames: []string{"hostname.com"},
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/rules", params.ProfileID), handler)
	actual, err := client.UpdateProfileCustomRule(context.Background(), params)

	group := int64(0)
	order := int64(1)
	want := []CustomRule{
		{
			Do:     Block,
			Status: IntBool{true},
			Group:  &group,
			Order:  &order,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfileCustomRule(context.Background(), UpdateProfileCustomRuleParams{
		ProfileID: "",
	})
	assert.Error(t, err, "Profile Custom Rule should not have been updated")
}

func TestDeleteProfileCustomRule(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": [],
			  "success": true,
			  "message": "Custom rule(s) deleted"
			}
		`)
	}

	params := DeleteProfileCustomRuleParams{
		ProfileID: "profileID",
		Hostname:  "hostname",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/rules/%s", params.ProfileID, params.Hostname), handler)
	actual, err := client.DeleteProfileCustomRule(context.Background(), params)

	want := []interface{}{}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.DeleteProfileCustomRule(context.Background(), DeleteProfileCustomRuleParams{
		ProfileID: "",
		Hostname:  "hostname",
	})
	assert.Error(t, err, "Profile Folder Custom Rule should not have been updated")

	_, err = client.DeleteProfileCustomRule(context.Background(), DeleteProfileCustomRuleParams{
		ProfileID: "profileID",
		Hostname:  "",
	})
	assert.Error(t, err, "Profile Custom Rule should not have been updated")
}
