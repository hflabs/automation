package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransitionToStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                string
		initialStatusId     string
		targetStatusId      string
		transitionsByStatus map[string][]Transition
		statusByTransition  map[string]string
		wantTransitions     []string
		wantStatusId        string
		wantErrContains     string
	}{
		{
			name:            "follows known intermediate statuses",
			initialStatusId: Issue.Status.New,
			targetStatusId:  Issue.Status.InProgressHRP,
			transitionsByStatus: map[string][]Transition{
				Issue.Status.New: {
					{ID: "close", Name: "Close", To: IssueField{ID: Issue.Status.NoNeedReaction}},
					{ID: "assign", Name: "Assign", To: IssueField{ID: Issue.Status.Assigned}},
				},
				Issue.Status.Assigned: {
					{ID: "start", Name: "Start", To: IssueField{ID: Issue.Status.InProgressHRP}},
				},
			},
			statusByTransition: map[string]string{
				"assign": Issue.Status.Assigned,
				"start":  Issue.Status.InProgressHRP,
			},
			wantTransitions: []string{"assign", "start"},
			wantStatusId:    Issue.Status.InProgressHRP,
		},
		{
			name:            "returns error when route is unknown",
			initialStatusId: Issue.Status.New,
			targetStatusId:  Issue.Status.InProgressHRP,
			transitionsByStatus: map[string][]Transition{
				Issue.Status.New: {
					{ID: "close", Name: "Close", To: IssueField{ID: Issue.Status.NoNeedReaction}},
				},
			},
			wantStatusId:    Issue.Status.New,
			wantErrContains: "cannot transition issue KEY-1",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			currentStatusId := tc.initialStatusId
			var appliedTransitions []string

			srv := newTransitionTestServer(t, &currentStatusId, &appliedTransitions, tc.transitionsByStatus, tc.statusByTransition)
			t.Cleanup(srv.Close)

			j := &jira{BaseUrl: srv.URL, Token: "token"}
			err := j.TransitionToStatus(ctx, "KEY-1", tc.targetStatusId)
			if tc.wantErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrContains)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantTransitions, appliedTransitions)
			require.Equal(t, tc.wantStatusId, currentStatusId)
		})
	}
}

func newTransitionTestServer(
	t *testing.T,
	currentStatusId *string,
	appliedTransitions *[]string,
	transitionsByStatus map[string][]Transition,
	statusByTransition map[string]string,
) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/issue/KEY-1", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, Issue.Fields.Status, r.URL.Query().Get("fields"))

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(IssueJira{
			Fields: FieldsIssue{Status: IssueField{ID: *currentStatusId}},
		})
		require.NoError(t, err)
	})
	mux.HandleFunc("/issue/KEY-1/transitions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			err := json.NewEncoder(w).Encode(TransitionsResponse{
				Transitions: transitionsByStatus[*currentStatusId],
			})
			require.NoError(t, err)
		case http.MethodPost:
			var req TransitionIssueRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)

			nextStatusId, ok := statusByTransition[req.Transition.ID]
			require.True(t, ok, "unexpected transition %q", req.Transition.ID)

			*appliedTransitions = append(*appliedTransitions, req.Transition.ID)
			*currentStatusId = nextStatusId
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	return httptest.NewServer(mux)
}
