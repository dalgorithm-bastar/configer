/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Using \"create\" command to acquire cfgfile locally",
    Long: `	create command works with 2 modes:
	offline: Totally generate locally
	online: Generate with local template and remote serviceinfo
 Add command "offline" or "online" to switch one`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("create called")
    },
}

func init() {
    rootCmd.AddCommand(createCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // createCmd.PersistentFlags().String("foo", "", "A help for foo")
    createCmd.PersistentFlags().StringVarP(&object.Version, "version", "v", "", "assign a config version(required)")
    createCmd.PersistentFlags().StringVarP(&object.Config, "config", "s", "", "assign config scheme(required)")
    createCmd.PersistentFlags().StringVarP(&object.Cluster, "cluster", "c", "", "assign a cluster name(required)")
    createCmd.PersistentFlags().StringVarP(&object.GlobalId, "globalid", "g", "", "assign a globalId(required)")
    createCmd.PersistentFlags().StringVarP(&object.LocalId, "localid", "l", "", "assign a localId within cluster(required)")
    createCmd.PersistentFlags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path, which default is pathin")
    createCmd.MarkPersistentFlagRequired("version")
    createCmd.MarkPersistentFlagRequired("env")
    createCmd.MarkPersistentFlagRequired("cluster")
    createCmd.MarkPersistentFlagRequired("globalid")
    createCmd.MarkPersistentFlagRequired("localid")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
