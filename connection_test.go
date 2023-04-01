package duty

import (
	"context"

	"github.com/flanksource/duty/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection", Ordered, func() {

	BeforeAll(func() {
		tx := testDB.Save(&models.Connection{
			Name:     "test",
			Type:     "test",
			Username: "configmap://test-cm/foo",
			Password: "secret://test-secret/foo",
			URL:      "sql://db?user=$(username)&password=$(password)",
		})
		Expect(tx.Error).ToNot(HaveOccurred())
	})

	It("username should be looked up from configmap", func() {

		user, err := GetEnvStringFromCache(testClient, "configmap://test-cm/foo", "default")
		Expect(err).ToNot(HaveOccurred())
		Expect(user).To(Equal("bar"))

		val, err := GetConfigMapFromCache(testClient, "default", "test-cm", "foo")
		Expect(err).ToNot(HaveOccurred())
		Expect(val).To(Equal("bar"))

	})
	var connection *models.Connection
	var err error
	It("should be retrieved successfully", func() {
		connection, err = GetConnection(context.Background(), testClient, testDB, "test", "test", "default")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should lookup kubernetes secrets", func() {
		Expect(connection.Username).To(Equal("bar"))
		Expect(connection.Password).To(Equal("secret"))
	})

	It("should template out the url", func() {
		Expect(connection.URL).To(Equal("sql://db?user=bar&password=secret"))
	})

})
