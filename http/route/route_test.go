// Integration tests.
// We could add "// +build integration" to the top of the file.
package route_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goyave.dev/goyave/v5/database"
	_ "goyave.dev/goyave/v5/database/dialect/postgres"
	"goyave.dev/goyave/v5/util/testutil"
	"goyave.dev/template/database/repository"
	"goyave.dev/template/http/dto"
	"goyave.dev/template/http/route"
	"goyave.dev/template/service/book"
	"goyave.dev/template/service/user"
)

func registerServices(server *testutil.TestServer) {
	userRepository := repository.NewUser(server.DB())
	userService := user.NewService(userRepository)
	server.RegisterService(userService)

	bookRepository := repository.NewBookSQL(server.DB())
	bookService := book.NewService(bookRepository)
	server.RegisterService(bookService)
}

func TestRoute(t *testing.T) {
	headerAuth := [2]string{"X-User-Id", "1"}

	type request struct {
		method  string
		url     string
		body    io.Reader
		headers [][2]string

		expectedStatus      int
		response            any
		quickBodyValidation string
	}

	// Could improve the testing set logic:
	// - with results injections:
	//   url: "/books?filter=id||$eq||{{ .Responses[0].id }}&per_page=1&page=1",
	// - with json schema validation: https://github.com/xeipuuv/gojsonschema
	testingSet := []struct {
		name     string
		requests []request
	}{
		{
			name: "I want to create a book and get my book ",
			requests: []request{
				{
					method:              "POST",
					url:                 "/books",
					body:                nil,
					headers:             [][2]string{headerAuth},
					expectedStatus:      201,
					response:            dto.Book{},
					quickBodyValidation: `{{ int .id | ne 0 }}`,
				},
				{
					method:              "GET",
					url:                 "/books?filter=id||$eq||{{ (index .Responses 0).id }}&per_page=1&page=1",
					body:                nil,
					headers:             [][2]string{headerAuth},
					expectedStatus:      200,
					response:            database.PaginatorDTO[dto.Book]{},
					quickBodyValidation: `{{ .records | len | eq 1 }}`,
				},
			},
		},
	}

	ts := testutil.NewTestServer(t, "config.test.json", nil)
	registerServices(ts)
	ts.RegisterRoutes(route.Register)

	for _, test := range testingSet {
		t.Run(test.name, func(t *testing.T) {
			tmplData := struct {
				Responses []any
			}{}
			tmpl := template.New("").Funcs(sprig.TxtFuncMap())

			for _, tr := range test.requests {
				{
					wr := bytes.NewBuffer(nil)
					template.Must(tmpl.Parse(tr.url)).Execute(wr, tmplData)
					tr.url = wr.String()
				}

				t.Logf("->: %s %s", tr.method, tr.url)

				req, err := http.NewRequest(tr.method, tr.url, tr.body)
				require.NoError(t, err)

				for _, header := range tr.headers {
					req.Header.Set(header[0], header[1])
				}

				resp := ts.TestRequest(req)
				require.NotNil(t, resp)

				var body any
				t.Logf("<-: %d", resp.StatusCode)
				func() {
					if resp.Body != nil {
						defer resp.Body.Close()
						bodyBytes, err := io.ReadAll(resp.Body)
						if !assert.NoError(t, err) {
							return
						}

						t.Logf("<-: %s", string(bodyBytes))

						require.NoError(t, json.Unmarshal(bodyBytes, &tr.response))
						body = tr.response

						if tr.quickBodyValidation == "" {
							return
						}

						wr := bytes.NewBuffer(nil)
						template.Must(tmpl.Parse(tr.quickBodyValidation)).Execute(wr, body)

						ok, _ := strconv.ParseBool(wr.String())
						t.Log(wr.String())
						require.True(t, ok, "body validation failed: %s -> %s", tr.quickBodyValidation, wr.String())
					}
				}()
				require.Equal(t, tr.expectedStatus, resp.StatusCode)

				tmplData.Responses = append(tmplData.Responses, body)
			}
		})
	}
}
