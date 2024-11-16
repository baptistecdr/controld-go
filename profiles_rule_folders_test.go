package controld

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListProfileRuleFolders(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"groups": [
				  {
					"PK": 1,
					"group": "Group Name",
					"action": {
					  "status": 1
					},
					"count": 0
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := ListProfileRuleFoldersParams{
		ProfileID: "profileID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/groups", params.ProfileID), handler)
	actual, err := client.ListProfileRuleFolders(context.Background(), params)

	want := []Group{
		{
			PK:    1,
			Group: "Group Name",
			Action: GroupAction{
				Do:     nil,
				Status: IntBool{true},
			},
			Count: 0,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.ListProfileRuleFolders(context.Background(), ListProfileRuleFoldersParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Rule Folders should not have been listed")
}

func TestCreateProfileRuleFolder(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"groups": [
				  {
					"PK": 1,
					"group": "New Group",
					"action": {
					  "status": 1
					},
					"count": 0
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := CreateProfileRuleFolderParams{
		ProfileID: "profileID",
		Name:      "New Group",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/groups", params.ProfileID), handler)
	actual, err := client.CreateProfileRuleFolder(context.Background(), params)

	want := []Group{
		{
			PK:    1,
			Group: "New Group",
			Action: GroupAction{
				Do:     nil,
				Status: IntBool{true},
			},
			Count: 0,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.CreateProfileRuleFolder(context.Background(), CreateProfileRuleFolderParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Rule FolderID should not have been created")
}

func TestUpdateProfileRuleFolder(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method, "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": {
				"groups": [
				  {
					"PK": 1,
					"group": "New Group Name",
					"action": {
					  "status": 1
					},
					"count": 0
				  }
				]
			  },
			  "success": true
			}
		`)
	}

	params := UpdateProfileRuleFolderParams{
		ProfileID: "profileID",
		FolderID:  "New Group Name",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/groups/%s", params.ProfileID, params.FolderID), handler)
	actual, err := client.UpdateProfileRuleFolder(context.Background(), params)

	want := []Group{
		{
			PK:    1,
			Group: "New Group Name",
			Action: GroupAction{
				Do:     nil,
				Status: IntBool{true},
			},
			Count: 0,
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.UpdateProfileRuleFolder(context.Background(), UpdateProfileRuleFolderParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Rule Folder should not have been updated")
}

func TestDeleteProfileRuleFolder(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method, "Expected method 'DELETE', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
			  "body": [],
			  "success": true,
			  "message": "Profile has been deleted"
			}
		`)
	}

	params := DeleteProfileRuleFolderParams{
		ProfileID: "profileID",
	}

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/groups/%s", params.ProfileID, params.FolderID), handler)
	actual, err := client.DeleteProfileRuleFolder(context.Background(), params)

	want := []interface{}{}
	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}

	_, err = client.DeleteProfileRuleFolder(context.Background(), DeleteProfileRuleFolderParams{
		ProfileID: "",
	})
	require.Error(t, err, "Profile Rule Folder should not have been updated")
}
