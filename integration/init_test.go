package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	buildpack               string
	cpythonBuildpack        string
	pipBuildpack            string
	pipInstallBuildpack     string
	minicondaBuildpack      string
	condaEnvUpdateBuildpack string

	buildpackInfo struct {
		Buildpack struct {
			ID   string
			Name string
		}
	}

	config struct {
		Cpython        string `json:"cpython"`
		Pip            string `json:"pip"`
		PipInstall     string `json:"pip-install"`
		Miniconda      string `json:"miniconda"`
		CondaEnvUpdate string `json:"conda-env-update"`
	}
)

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	root, err := filepath.Abs("./..")
	Expect(err).ToNot(HaveOccurred())

	file, err := os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&buildpackInfo)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	file, err = os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&config)).To(Succeed())
	Expect(file.Close()).To(Succeed())

	buildpackStore := occam.NewBuildpackStore()

	buildpack, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	cpythonBuildpack, err = buildpackStore.Get.
		Execute(config.Cpython)
	Expect(err).NotTo(HaveOccurred())

	pipBuildpack, err = buildpackStore.Get.
		Execute(config.Pip)
	Expect(err).NotTo(HaveOccurred())

	pipInstallBuildpack, err = buildpackStore.Get.
		Execute(config.PipInstall)
	Expect(err).NotTo(HaveOccurred())

	minicondaBuildpack, err = buildpackStore.Get.
		Execute(config.Miniconda)
	Expect(err).NotTo(HaveOccurred())

	condaEnvUpdateBuildpack, err = buildpackStore.Get.
		Execute(config.CondaEnvUpdate)
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(5 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Default", testDefault)
	suite.Run(t)
}
