package main_test

import (
	"code.cloudfoundry.org/cli/util/testhelpers/pluginbuilder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWhoamiPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	pluginbuilder.BuildTestBinary(".", "main")
	RunSpecs(t, "TellMeAJokePlugin Suite")
}
