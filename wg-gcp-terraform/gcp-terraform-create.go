package wg-gcp-terraform

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	// https://github.com/hashicorp/terraform-exec
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func wg_gcp_terraform() {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpDir, false))
	if err != nil {
		panic(err)
	}

	workingDir := "/home/Stanislav/wireguard-bot/gcp-terraform/infrastructure"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		panic(err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	if err != nil {
		panic(err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(state.FormatVersion) // "0.1"
}
