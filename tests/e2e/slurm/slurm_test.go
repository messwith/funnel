package slurm

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/tests/e2e"
	"os"
	"testing"
)

var fun *e2e.Funnel

func TestMain(m *testing.M) {
	conf := e2e.DefaultConfig()
	if conf.Backend != "slurm" {
		logger.Debug("Skipping slurm e2e tests...")
		os.Exit(0)
	}

	fun = e2e.NewFunnel(conf)
	fun.StartServerInDocker("ohsucompbio/slurm:latest", []string{"--hostname", "ernie"})
	defer fun.CleanupTestServerContainer()

	m.Run()
	return
}

func TestHelloWorld(t *testing.T) {
	id := fun.Run(`
    --sh 'echo hello world'
  `)
	task := fun.Wait(id)

	if task.Logs[0].Logs[0].Stdout != "hello world\n" {
		t.Fatal("Missing stdout")
	}
}

func TestResourceRequest(t *testing.T) {
	id := fun.Run(`
    --sh 'echo I need resources!' --cpu 1 --ram 2 --disk 5
  `)
	task := fun.Wait(id)

	if task.Logs[0].Logs[0].Stdout != "I need resources!\n" {
		t.Fatal("Missing stdout")
	}
}
