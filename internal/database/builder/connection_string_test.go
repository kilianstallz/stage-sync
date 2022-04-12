package builder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"stage-sync-cli/config"
	"stage-sync-cli/internal/database/builder"
	"testing"
)

func TestBuildConnectionString(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BuildConnectionString Suite")
}

var _ = Describe("Build Connection Strings", func () {
	Describe("Postgres", func () {
		It("should build a connection string", func () {
			config := config.ConfigDB{
				User:     "user",
				Password: "password",
				Host:     "host",
				Port:     5432,
				Database: "database",
				SslMode: "require",
			}
			result := builder.BuildConnectionString(config)
			Expect(result).To(Equal("postgres://user:password@host:5432/database?sslmode=require"))
		})
		It("should build without sslmode", func () {
			config := config.ConfigDB{
				User:     "user",
				Password: "password",
				Host:     "host",
				Port:     5432,
				Database: "database",
			}
			result := builder.BuildConnectionString(config)
			Expect(result).To(Equal("postgres://user:password@host:5432/database"))
		})
		It("should build without creds", func () {
			config := config.ConfigDB{
				Host:     "host",
				Port:     5432,
				Database: "database",
			}
			result := builder.BuildConnectionString(config)
			Expect(result).To(Equal("postgres://host:5432/database"))
		})
	})
})


