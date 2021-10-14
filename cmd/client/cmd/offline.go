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
	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/template"
	"github.com/configcenter/pkg/util"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// offlineCmd represents the offline command
var offlineCmd = &cobra.Command{
	Use:   "offline",
	Short: "Generate offline",
	Long:  `Generate configfile under mode offline`,
	Example: `$ ./cfgtool create offline --version 1.0.0 --env 00 --cluster EzEI.set0 --globalid 141 --localid 3 --pathin /etc/configcenter/configure.tar.gz --pathout /etc/configcenter
$ ./cfgtool create offline -v 1.0.0 -e 00 -c EzEI.set0 -g 141 -l 3 -i /etc/configcenter/configure.tar.gz -o /etc/configcenter`,
	Run: CreateOffline,
}

func init() {
	createCmd.AddCommand(offlineCmd)
	offlineCmd.Flags().StringVarP(&object.PathIn, "pathin", "i", "", "assign input path of compressedfile(required)")
	offlineCmd.Flags().StringVarP(&object.TemplateName, "template", "t", "", "assign template to fill(required)")
	offlineCmd.MarkFlagRequired("template")
	offlineCmd.MarkFlagRequired("pathin")
}

func CreateOffline(cmd *cobra.Command, args []string) {
	//获取文件存储路径和内容
	fileName := filepath.Base(object.PathIn)
	if filepath.Ext(fileName) != ".zip" && !strings.Contains(fileName, ".tar.gz") {
		fmt.Println(fmt.Sprintf("unsupport compressed file of \"%s\"", fileName))
	}
	//TODO 文件夹数据
	//初始化repository，客户端为压缩包模式
	err := repository.NewStorage(repository.CompressedFileType, object.PathIn)
	if err != nil {
		fmt.Println(err)
		return
	}
	path := util.Join("/", object.Version, object.Env, object.Cluster, repository.Templates, object.TemplateName)
	tmplContent, err := repository.Src.Get(path)
	test := repository.Src.GetSourceDataorOperator()
	fmt.Sprint(test)
	if err != nil {
		fmt.Println(err)
		return
	}
	templateIns, err := template.NewTemplateImpl(repository.Src, object.GlobalId, object.LocalId, "tmplIns", object.Version, object.Env)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := templateIns.Fill(tmplContent, object.TemplateName)
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建文件夹和文件
	dirPath := filepath.Dir(object.PathIn)
	if object.PathOut != "" {
		dirPath = object.PathOut
	}
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	outPutFileName := util.Join("_", object.Version, object.Env, object.Cluster, object.LocalId, object.TemplateName)
	f, err := os.OpenFile(dirPath+"/"+outPutFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
		fmt.Println(err)
		return
	}
}
