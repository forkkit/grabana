package grabana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoldersCanBeCreated(t *testing.T) {
	req := require.New(t)
	folderName := "Test folder"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "uid": "nErXDvCkzz",
  "id": 1,
  "title": "Test folder"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	folder, err := client.CreateFolder(context.TODO(), folderName)

	req.NoError(err)
	req.Equal(folderName, folder.Title)
}

func TestFoldersCreationCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{
  "message": "The folder has been changed by someone else",
  "status": "version-mismatch"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	folder, err := client.CreateFolder(context.TODO(), "")

	req.Error(err)
	req.Nil(folder)
}

func TestAFolderCanBeFoundByTitle(t *testing.T) {
	req := require.New(t)
	folderName := "Test folder"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "nErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "nErXDvCkyy",
    "title": "Test folder"
  },
  {
    "id":3,
    "uid": "nErXDvCkxx",
    "title": "Department XYZ"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	folder, err := client.GetFolderByTitle(context.TODO(), strings.ToLower(folderName))

	req.NoError(err)
	req.Equal(folderName, folder.Title)
}

func TestAnExplicitErrorIsReturnedIfTheFolderIsNotFound(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
  {
    "id":1,
    "uid": "nErXDvCkzz",
    "title": "Department ABC"
  },
  {
    "id":2,
    "uid": "nErXDvCkyy",
    "title": "Test folder"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	folder, err := client.GetFolderByTitle(context.TODO(), "folder that do not exist")

	req.Error(err)
	req.Nil(folder)
	req.Equal(ErrFolderNotFound, err)
}

func TestGetFolderByTitleCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	folder, err := client.GetFolderByTitle(context.TODO(), "folder that do not exist")

	req.Error(err)
	req.Nil(folder)
}

func TestAnAlertChannelCanBeFoundByName(t *testing.T) {
	req := require.New(t)
	name := "Team B"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
  {
    "id": 1,
    "uid": "team-a-email-notifier",
    "name": "Team A",
    "type": "email"
  },
  {
    "id": 2,
    "uid": "team-b-email-notifier",
    "name": "Team B",
    "type": "email"
  },
  {
    "id": 1,
    "uid": "team-c-email-notifier",
    "name": "Team C",
    "type": "email"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	channel, err := client.GetAlertChannelByName(context.TODO(), strings.ToLower(name))

	req.NoError(err)
	req.Equal(name, channel.Name)
	req.Equal("email", channel.Type)
}

func TestAnExplicitErrorIsReturnedIfTheChannelIsNotFound(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
   {
    "id": 1,
    "uid": "team-a-email-notifier",
    "name": "Team A",
    "type": "email"
  },
  {
    "id": 2,
    "uid": "team-b-email-notifier",
    "name": "Team B",
    "type": "email"
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	channel, err := client.GetAlertChannelByName(context.TODO(), "channel that do not exist")

	req.Error(err)
	req.Nil(channel)
	req.Equal(ErrAlertChannelNotFound, err)
}

func TestGetAlertChannelByNameCanFail(t *testing.T) {
	req := require.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{}}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	folder, err := client.GetAlertChannelByName(context.TODO(), "channel that do not exist")

	req.Error(err)
	req.Nil(folder)
}

func TestDashboardsCanBeCreated(t *testing.T) {
	req := require.New(t)
	dashboard := NewDashboardBuilder("Dashboard name")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "id":      1,
  "uid":     "cIBgcSjkk",
  "url":     "/d/cIBgcSjkk/production-overview",
  "status":  "success",
  "version": 1,
  "slug":    "production-overview"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	board, err := client.UpsertDashboard(context.TODO(), &Folder{}, dashboard)

	req.NoError(err)
	req.NotNil(board)
}

func TestDashboardsCreationCanFail(t *testing.T) {
	req := require.New(t)
	dashboard := NewDashboardBuilder("Dashboard name")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{
  "message": "The folder has been changed by someone else",
  "status": "version-mismatch"
}`)
	}))
	defer ts.Close()

	client := NewClient(http.DefaultClient, ts.URL, "")

	board, err := client.UpsertDashboard(context.TODO(), &Folder{}, dashboard)

	req.Error(err)
	req.Nil(board)
}
