package controller_test

import (
	"fmt"
	"log"
	"os"

	"github.com/midoblgsm/ubiquity/utils/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var testLogger *log.Logger
var logFile *os.File

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	defer logs.InitStdoutLogger(logs.DEBUG)()

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeEach(func() {
	var err error
	logFile, err = os.OpenFile("/tmp/test-ubiquity-csi.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to setup logger: %s\n", err.Error())
		return
	}
	testLogger = log.New(logFile, "ubiquity-csi: ", log.Lshortfile|log.LstdFlags)
})

var _ = AfterEach(func() {
	err := logFile.Sync()
	if err != nil {
		panic(err.Error())
	}
	err = logFile.Close()
	if err != nil {
		panic(err.Error())
	}
})
