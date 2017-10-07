package main_test

import (
	"net/http"
	"strings"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	io_helpers "code.cloudfoundry.org/cli/util/testhelpers/io"

	. "github.com/jtuchscherer/tell-me-a-joke-plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("TellMeAJokePlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *pluginfakes.FakeCliConnection
		var tellMeAJokeCmd *TellMeAJokeCmd
		var outputChan chan []string
		var server *ghttp.Server
		var statusCode int

		BeforeEach(func() {
			statusCode = http.StatusOK
			outputChan = make(chan []string)

			fakeCliConnection = &pluginfakes.FakeCliConnection{}
			fakeCliConnection.UsernameReturns("user@user.com", nil)
			server = ghttp.NewServer()

			tellMeAJokeCmd = NewTellMeAJokeCmd(server.URL())

		})

		Context("when there is internet connection to the joke api", func() {
			
				Context("when the response is proper JSON", func() {
					BeforeEach(func() {

						server.AppendHandlers(ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/"),
							ghttp.VerifyHeaderKV("Accept", "application/json"),
							ghttp.VerifyHeaderKV("User-Agent", "https://github.com/jtuchscherer/tell-me-joke-plugin"),
							ghttp.RespondWith(200, `
        				{
                  "id": "R7UfaahVfFd",
                  "joke": "My dog used to chase people on a bike a lot. It got so bad I had to take his bike away.",
                  "status": 200
                }
				      `),
						))
					})
					It("displays the joke's text", func(done Done) {
						defer close(done)
						go invokeCmd(outputChan, tellMeAJokeCmd, fakeCliConnection)

						var output []string
						Eventually(outputChan, 2).Should(Receive(&output))
						outputString := strings.Join(output, "")
						Expect(outputString).To(Equal("My dog used to chase people on a bike a lot. It got so bad I had to take his bike away."))
					})
				})

				Context("when the response is not in the expected JSON", func() {
					BeforeEach(func() {
						server.AppendHandlers(ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/"),
							ghttp.VerifyHeaderKV("Accept", "application/json"),
							ghttp.VerifyHeaderKV("User-Agent", "https://github.com/jtuchscherer/tell-me-joke-plugin"),
							ghttp.RespondWith(200, `something else`),
						))
					})
					It("shows a sad error message", func(done Done) {
						defer close(done)
						go invokeCmd(outputChan, tellMeAJokeCmd, fakeCliConnection)

						var output []string
						Eventually(outputChan, 2).Should(Receive(&output))
						outputString := strings.Join(output, "")
						Expect(outputString).To(ContainSubstring("FAILED"))
					})
				})
					Context("when the response is empty", func() {
					BeforeEach(func() {
						server.AppendHandlers(ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/"),
							ghttp.VerifyHeaderKV("Accept", "application/json"),
							ghttp.VerifyHeaderKV("User-Agent", "https://github.com/jtuchscherer/tell-me-joke-plugin"),
							ghttp.RespondWith(200, ""),
						))
					})
					It("shows a sad error message", func(done Done) {
						defer close(done)
						go invokeCmd(outputChan, tellMeAJokeCmd, fakeCliConnection)

						var output []string
						Eventually(outputChan, 2).Should(Receive(&output))
						outputString := strings.Join(output, "")
						Expect(outputString).To(ContainSubstring("FAILED"))
					})
				})
				
			

		})
		Context("when there is no connection to the joke api", func() {
			BeforeEach(func() {

				server.Close()
			})
			It("shows a good error message", func(done Done) {
				defer close(done)
				go invokeCmd(outputChan, tellMeAJokeCmd, fakeCliConnection)

				var output []string
				Eventually(outputChan, 2).Should(Receive(&output))
				outputString := strings.Join(output, "")
				Expect(outputString).To(ContainSubstring("FAILED"))
			})
		})

	})
})

func invokeCmd(outputChan chan []string, tellMeAJokeCmd *TellMeAJokeCmd, fakeCliConnection *pluginfakes.FakeCliConnection) {
	outputChan <- io_helpers.CaptureOutput(func() {
		tellMeAJokeCmd.Run(fakeCliConnection, []string{"tell-me-a-joke"})
	})
}
