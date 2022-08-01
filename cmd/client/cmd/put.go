package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/configcenter/pkg/define"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/util"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "put compressed configfile to remote",
	Long: `put command enables you to save config file under your username on remote
attention that the configfile would not be submit.
learn more about that on command "post"`,
	Run: Put,
}

func init() {
	rootCmd.AddCommand(putCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	putCmd.Flags().StringVarP(&object.PathIn, "pathin", "i", "", "assign input path of compressedfile(required)")
	putCmd.MarkFlagRequired("pathin")
}

func Put(cmd *cobra.Command, args []string) {
	//检测是否为文件夹
	object.PathIn = filepath.Clean(object.PathIn)
	s, err := os.Stat(object.PathIn)
	if err != nil || !s.IsDir() {
		fmt.Println("please specify input path to directory")
		return
	}
	/*fileMap := make(map[string][]byte)
	  //获取权限文件
	  permData, err := generatePermFile()
	  if err != nil {
	      fmt.Println(err)
	      return
	  }
	  pathSli := strings.Split(object.PathIn, _separator)
	  //插入整个版本对应的权限信息文件
	  fileMap[pathSli[len(pathSli)-1]+"/"+define.Perms] = permData
	  lenVersion := len(pathSli[len(pathSli)-1])
	  idx := len(object.PathIn) - lenVersion + 1
	  err = filepath.Walk(object.PathIn, func(path string, info os.FileInfo, err error) error {
	      if err != nil {
	          fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
	          return err
	      }
	      if !info.IsDir() {
	          pathKey := path[idx-1:]
	          data, readErr := ioutil.ReadFile(path)
	          if readErr != nil {
	              return readErr
	          }
	          //转化为标准路径分隔符
	          fileMap[filepath.ToSlash(pathKey)] = data
	      }
	      return err
	  })
	  if err != nil {
	      fmt.Printf("error when getting permission value walking the path %q: %v\n", object.PathIn, err)
	      return
	  }*/
	fileMap, _, err := util.LoadDirWithPermFile(object.PathIn, util.Separator, define.Template)
	compressedFileData, err := util.CompressToStream("cfgpkg.tar.gz", fileMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = checkInputPackage(fileMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	//新建客户端
	//读取grpc配置信息
	err = GetGrpcClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	//新建grpc客户端
	conn, err := grpc.Dial(GrpcInfo.Socket, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	client := pb.NewConfigCenterClient(conn)
	resp, err := client.PUT(context.Background(), &pb.CfgReq{
		UserName: object.UserName,
		File: &pb.AnyFile{
			FileName: "cfgpkg.tar.gz",
			FileData: compressedFileData,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.Status != "ok" {
		fmt.Println("resp err: " + resp.Status)
		return
	}
	fmt.Println(fmt.Sprintf("Put cfgpkg to remote succeed"))
}

// 进行标准比较，无需考虑系统路径分隔符
func checkInputPackage(fileMap map[string][]byte) error {
	for filePath, _ := range fileMap {
		if strings.Contains(filePath, define.Deployment) {
			pathSli := strings.SplitN(filePath, "/", 7)
			if len(pathSli) != 7 {
				return fmt.Errorf("err deployment file path of:%s, please checkout input path or pkg format", filePath)
			}
			if pathSli[4] != define.DeploymentFlag {
				return fmt.Errorf("err deployment file path of:%s, differ from standard path with flag: %s", filePath, define.DeploymentFlag)
			}
		} else if strings.Contains(filePath, define.Service) {
			pathSli := strings.SplitN(filePath, "/", 6)
			if len(pathSli) < 6 {
				return fmt.Errorf("err service file path of:%s, please checkout input path or pkg format", filePath)
			}
			if pathSli[4] != define.ServiceFlag {
				return fmt.Errorf("err service file path of:%s, differ from standard path with flag: %s", filePath, define.ServiceFlag)
			}
		} else if strings.Contains(filePath, define.Template) {
			pathSli := strings.SplitN(filePath, "/", 6)
			if len(pathSli) < 5 || pathSli[4] != define.TemplateFlag {
				return fmt.Errorf("err template file path of:%s, please checkout input path or pkg format", filePath)
			}
		} else {
			//处理剩余情况
			if !strings.Contains(filePath, define.Infrastructure) && !strings.Contains(filePath, define.Perms) {
				return fmt.Errorf("a file should not be here with filepath: %s", filePath)
			}
		}
	}
	return nil
}
