package config_test

import (
	. "github.com/onsi/gomega"
	"stage-sync/internal/config"
	"testing"
)

func Test_ConfigSuite(t *testing.T) {
	RegisterTestingT(t)
	// should return a valid configuration
	cnfg, err := config.ParseConfigFromFile("./mocks/valid.yaml")
	Expect(err).To(BeNil())
	Expect(cnfg.SourceDatabase).To(Not(BeNil()))
	Expect(cnfg.TargetDatabase).To(Not(BeNil()))
	Expect(cnfg.Tables).To(Not(BeNil()))
	Expect(len(cnfg.Tables)).ShouldNot(BeZero())
	Expect(len(cnfg.Tables)).Should(Equal(2))
	Expect(len(cnfg.Tables[0].PrimaryKeys)).Should(Equal(1))

	Expect(cnfg.Tables[0].PrimaryKeys).Should(ContainElement("Id"))
	Expect(cnfg.Tables[1].PrimaryKeys).Should(ContainElement("Id"))

	Expect(len(cnfg.Tables[1].OnlyWhere)).Should(Equal(1))

	Expect(cnfg.Tables[0].NoDelete).To(BeFalse())
	Expect(cnfg.Tables[1].NoDelete).To(BeTrue())

	// should fail on missing file
	_, err = config.ParseConfigFromFile("./mocks/missing.yaml")
	Expect(err).ToNot(BeNil())

	// should fail on invalid file
	_, err = config.ParseConfigFromFile("./mocks/invalid.js")
	Expect(err).ToNot(BeNil())
}
