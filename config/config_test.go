package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"stage-sync-cli/config"
	"testing"
)

func Test_ConfigSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Parse Configuration file", func() {
	Context("when the configuration file is valid", func() {
		It("should return a valid configuration", func() {
			config, err := config.ParseConfigFromFile("./mocks/valid.yaml")
			Expect(err).To(BeNil())
			Expect(config.SourceDatabase).To(Not(BeNil()))
			Expect(config.TargetDatabase).To(Not(BeNil()))
			Expect(config.Tables).To(Not(BeNil()))
			Expect(len(config.Tables)).ShouldNot(BeZero())
			Expect(len(config.Tables)).Should(Equal(2))
			Expect(len(config.Tables[0].PrimaryKeys)).Should(Equal(1))

			Expect(config.Tables[0].PrimaryKeys).Should(ContainElement("Id"))
			Expect(config.Tables[1].PrimaryKeys).Should(ContainElement("Id"))

			Expect(len(config.Tables[1].OnlyWhere)).Should(Equal(1))

			Expect(config.Tables[0].NoDelete).To(BeFalse())
			Expect(config.Tables[1].NoDelete).To(BeTrue())
		})

		It("should fail on missing file", func() {
			_, err := config.ParseConfigFromFile("./mocks/missing.yaml")
			Expect(err).ToNot(BeNil())
		})

		It("should fail on invalid file", func() {
			_, err := config.ParseConfigFromFile("./type.go")
			Expect(err).ToNot(BeNil())
		})
	})
})
