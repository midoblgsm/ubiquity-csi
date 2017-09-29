# Ubiquity CSI Plugin
[![Build Status](https://travis-ci.org/midoblgsm/ubiquity-csi.svg?branch=master)](https://travis-ci.org/midoblgsm/ubiquity-csi)
[![GoDoc](https://godoc.org/github.com/midoblgsm/ubiquity-csi?status.svg)](https://godoc.org/github.com/midoblgsm/ubiquity-csi)
[![Coverage Status](https://coveralls.io/repos/github/midoblgsm/ubiquity-csi/badge.svg?branch=master)](https://coveralls.io/github/midoblgsm/ubiquity-csi?branch=master)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/midoblgsm/ubiquity-csi)](https://goreportcard.com/report/github.com/midoblgsm/ubiquity-csi)






This project includes a [CSI](https://github.com/container-storage-interface/spec) plugin for managing persistent volumes using [Ubiquity](https://github.com/IBM/ubiquity) service.
The code is provided as is, without warranty. Any issue will be handled on a best-effort basis.

### Prerequisites
  * Install [golang](https://golang.org/) (>=1.6).
  * Install [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git).
  * Configure go. GOPATH environment variable must be set correctly before starting the build process. Create a new directory and set it as GOPATH.
  * Setup [Ubiquity Service](https://github.com/midoblgsm/ubiquity)
### Download and build source code
* Build Ubiquity service from source. 
```bash
mkdir -p $HOME/workspace
export GOPATH=$HOME/workspace
mkdir -p $GOPATH/src/github.com/midoblgsm
cd $GOPATH/src/github.com/midoblgsm
git clone git@github.com:midoblgsm/ubiquity-csi.git
cd ubiquity-csi
./scripts/run_glide_up
./scripts/build
```
The built binary will be in the bin directory inside the repository folder.
To run it you need to setup ubiquity client configuration in [ubiquity-client.conf](ubiquity-client.conf).
After that setup the CSI endpoint:
```bash
export CSI_ENDPOINT=tcp://127.0.0.1:9595
./bin/ubiquity-csi
```

### Running unit tests for ubiquity-csi

Install these go packages to test Ubiquity:
```bash
# Install ginkgo
go install github.com/onsi/ginkgo/ginkgo
# Install gomega
go install github.com/onsi/gomega
```

Run the tests:
```bash
./scripts/run_units.sh
```


### Test using CSI client
The CSI client is supposed to mimic the Container frameworks calls. 
In order to build the client run the following commands inside the repo directory:
```bash
./script/build_client
```
The binary will be built inside bin.

Then you can use it to execute commands like:

```bash
# create a volume named testVolume
./bin/ubiquity-csi-client createvolume -endpoint tcp://127.0.0.1:9595 -limitBytes 512 -o nfs -requiredBytes 512 -service gold -t xfs -version 0.0.0 -params {\"backend\":\"localhost\"} testVolume
# or List the existing volumes
./bin/ubiquity-csi-client listvolumes -endpoint tcp://127.0.0.1:9595

```            
           
### Support
For any questions, suggestions, or issues, use github.