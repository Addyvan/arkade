// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallArgo() *cobra.Command {
	var command = &cobra.Command{
		Use:          "argo",
		Short:        "Install argo",
		Long:         `Install argo`,
		Example:      `  arkade install argo`,
		SilenceUsage: true,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		_, err := k8s.KubectlTask("create", "ns",
			"argo")
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}

		_, err = k8s.KubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/argoproj/argo-workflows/v2.11.8/manifests/install.yaml", "-n", "argo")
		if err != nil {
			return err
		}

		fmt.Println(ArgoInfoMsgInstallMsg)

		return nil
	}

	return command
}

const ArgoInfoMsg = `
# Get the Argo CLI
arkade install argo

# Port-forward the Argo API server
kubectl -n argo port-forward deployment/argo-server 2746:2746

# Open the UI:
https://localhost:2746

# Get started with Argo Workflows at
# https://argoproj.github.io/argo-workflows/quick-start/`

const ArgoInfoMsgInstallMsg = `=======================================================================
= Argo has been installed                                           =
=======================================================================` +
	"\n\n" + ArgoInfoMsg + "\n\n" + pkg.ThanksForUsing
