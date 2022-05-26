package graph

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	gqlclient "github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	c "github.com/Littlefisher619/cosdisk/config"
	"github.com/Littlefisher619/cosdisk/graph/auth"
	"github.com/Littlefisher619/cosdisk/graph/generated"
	"github.com/Littlefisher619/cosdisk/graph/model"
	m "github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/repository/dbdriver"
	"github.com/Littlefisher619/cosdisk/service"
	"github.com/stretchr/testify/require"
)

func authTestDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return next(ctx)
}

func TestFile(t *testing.T) {
	models := []interface{}{
		&m.User{},
		&m.ShareFile{},
	}
	db, err := dbdriver.CreateXormEngine("mysql", "root@/tcp(127.0.0.1:4000)", models, true)
	if err != nil {
		panic(err)
	}
	defer db.Db.Close()
	rc := dbdriver.InitTikv([]string{"127.0.0.1:2379"})
	if rc == nil {
		return
	}
	defer rc.CloseTikv()

	c, err := c.LoadConfig("config.toml")
	if err != nil {
		panic(err)
	}
	service := service.New(c)
	/*
		service := service.New(
			account.NewMap(),
			userfile.NewMap(),
			sharefile.New(),
		)
	*/
	jwtAuthManager := auth.New()
	graphqlConfig := generated.Config{Resolvers: &Resolver{service, jwtAuthManager}}
	graphqlConfig.Directives.Auth = authTestDirective
	srv := httptest.NewServer(handler.NewDefaultServer(generated.NewExecutableSchema(graphqlConfig)))
	defer srv.Close()

	gql := gqlclient.New(srv.Config.Handler, gqlclient.Path("/graphql"))

	aTxtFile, _ := ioutil.TempFile(os.TempDir(), "a.txt")
	defer os.Remove(aTxtFile.Name())
	aTxtFile.WriteString(`test`)

	t.Run("register user", func(t *testing.T) {
		var result interface{}

		mutation := `mutation {
			register(input: { name: "fish", email: "fish", password: "123" }) {
			  token
			  user {
				id
				name
				email
			  }
			}
		  }`

		err := gql.Post(mutation, &result)
		require.Nil(t, err)
	})

	t.Run("login user", func(t *testing.T) {
		var result interface{}

		mutation := `mutation {
			login(email: "fish", password: "123") {
			  token
			  user {
				id
				name
				email
			  }
			}
		  }`

		err := gql.Post(mutation, &result)
		require.Nil(t, err)

		// test wrong password
		mutation = `mutation {
			login(email: "fish", password: "fff") {
			  token
			  user {
				id
				name
				email
			  }
			}
		  }`

		err = gql.Post(mutation, &result)
		require.Error(t, err)
	})

	t.Run("single file upload", func(t *testing.T) {
		var resString interface{}

		mutation := `mutation {
			createDir(path: "/")
		  }`

		err := gql.Post(mutation, &resString)
		require.Nil(t, err)

		mutation = `mutation {
			createDir(path: "/testdir")
		  }`

		err = gql.Post(mutation, &resString)
		require.Nil(t, err)

		mutation = `mutation ($file: Upload!) {
			singleUpload(input: { path: "/a.txt", file: $file }) {
				name
				contentType
			}
		}`
		var result struct {
			SingleUpload *model.File
		}

		err = gql.Post(mutation, &result, gqlclient.Var("file", aTxtFile), gqlclient.WithFiles())
		require.Nil(t, err)
		require.Contains(t, result.SingleUpload.Name, "a.txt")
		require.Equal(t, "text/plain; charset=utf-8", result.SingleUpload.ContentType)
	})

	t.Run("single file download url", func(t *testing.T) {
		// fix me: require secret to download from cos
		var result interface{}

		mutation := `mutation {
			getDownloadURL(path: "/a.txt")
		  }`

		err := gql.Post(mutation, &result)
		require.Nil(t, err)
	})

	t.Run("list dir", func(t *testing.T) {
		var result struct {
			ListDir []model.File `json:"listDir"`
		}
		mutation := `mutation {
			listDir(path: "/") {
				  name
				  contentType
			}
		  }`

		err := gql.Post(mutation, &result)
		require.Nil(t, err)
		require.EqualValues(t, []model.File{{Name: "testdir", ContentType: "dir"},
			{Name: "a.txt", ContentType: "file"}}, result.ListDir)
	})

	t.Run("delete file and dir", func(t *testing.T) {
		var resString1 interface{}
		mutation := `mutation {
			createDir(path: "/dir")
		}`
		err := gql.Post(mutation, &resString1)
		require.Nil(t, err)

		var resString2 interface{}
		mutation = `mutation {
			deleteFile(path: "/a.txt")
		  }`
		err = gql.Post(mutation, &resString2)
		require.Nil(t, err)

		var resString3 interface{}
		mutation = `mutation {
			deleteDir(path: "/dir")
		  }`
		err = gql.Post(mutation, &resString3)
		require.Nil(t, err)
	})
}

// test curl:
//curl -H "Content-Type:application/json" -X POST  --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImV4cCI6MTY0NDM4NTQ0MSwiaWF0IjoxNjQ0MTI2MjQxfQ.hTlGa4I5edA1-klC19n9WT6QQ3EHKPkfrAn9_xF13pI" --data '{"query":"mutation { createDir(path: \"/\")}","variables":{},"operationName":null}' http://127.0.0.1:8080/query
//curl -H "Content-Type:application/json" -X POST --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImV4cCI6MTY0NDI1MzQ2MSwiaWF0IjoxNjQzOTk0MjYxfQ._FTiN_vhu2K-Wl9j1Gi6ylOdu1IkDwm-aaEdWCy7PGw" --data '{"query":"\n\nmutation {\n  auth {\n    register(input: { name: \"fish\", email: \"string\", password: \"123\" })\n  }\n}","variables":{},"operationName":null}' http://127.0.0.1:8080/query
