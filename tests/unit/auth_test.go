package unit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"morellis/app"
	"morellis/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Load JWT Auth test file")
}

var _ = Describe("JWT Auth handles request...", func() {
	Context("when path is in notAuth...", func() {
		It("should serve the request...", func() {
			ts := httptest.NewServer(app.JwtAuthentication(GetTestHandler()))
			defer ts.Close()

			var u bytes.Buffer
			u.WriteString(string(ts.URL))
			u.WriteString("/api/user/new")

			res, err := http.Get(u.String())

			Expect(err).To(BeNil())
			if res != nil {
				defer res.Body.Close()
			}
			_, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(BeEquivalentTo(200))
		})
	})
	Context("when path is not in notAuth", func() {
		It("should deny requests without a token", func() {
			ts := httptest.NewServer(app.JwtAuthentication(GetTestHandler()))
			defer ts.Close()

			var u bytes.Buffer
			u.WriteString(string(ts.URL))
			u.WriteString("/api/flavor")

			res, err := http.Get(u.String())

			Expect(err).To(BeNil())
			if res != nil {
				defer res.Body.Close()
			}
			_, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(BeEquivalentTo(403))
		})

		It("should serve requests with a valid token", func() {
			ts := httptest.NewServer(app.JwtAuthentication(GetTestHandler()))
			defer ts.Close()

			var u bytes.Buffer
			u.WriteString(string(ts.URL))
			u.WriteString("/api/flavor")

			tk := &models.Token{UserId: 1}
			token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
			tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

			req, err := http.NewRequest("GET", u.String(), nil)
			Expect(err).To(BeNil())
			req.Header.Add("authorization", fmt.Sprintf("Bearer %v", tokenString))

			client := http.Client{}
			res, err := client.Do(req)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(BeEquivalentTo(200))
		})
	})
})
