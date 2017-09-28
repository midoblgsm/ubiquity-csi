/**
 * Copyright 2017 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	ctl "github.com/midoblgsm/ubiquity-csi/controller"
	"github.com/midoblgsm/ubiquity/fakes"
	"github.com/midoblgsm/ubiquity/resources"
)

var _ = Describe("Controller", func() {

	var (
		fakeClient     *fakes.FakeStorageClient
		controller     *ctl.Controller
		fakeExec       *fakes.FakeExecutor
		ubiquityConfig resources.UbiquityPluginConfig
	)
	BeforeEach(func() {
		fakeExec = new(fakes.FakeExecutor)
		ubiquityConfig = resources.UbiquityPluginConfig{}
		fakeClient = new(fakes.FakeStorageClient)
		controller = ctl.NewControllerWithClient(testLogger, fakeClient, fakeExec)
		os.MkdirAll("/tmp/test/mnt2", 0777)
	})

	Context(".CreateVolume", func() {
		It("Should fail when ubiquity client returns an error", func() {
			params := make(map[string]string)
			params["backend"] = "test_backend"
			createResponse := resources.CreateVolumeResponse{Volume: resources.Volume{}, Error: fmt.Errorf("error occurred")}
			fakeClient.CreateVolumeReturns(createResponse)

			request := csi.CreateVolumeRequest{Name: "testVolume", Version: &csi.Version{}, CapacityRange: &csi.CapacityRange{RequiredBytes: 100, LimitBytes: 100}, Parameters: params}

			createVolumeResponse, err := controller.CreateVolume(request)
			Expect(err).To(HaveOccurred())
			Expect(createVolumeResponse).ToNot(BeNil())
		})
	})
})
