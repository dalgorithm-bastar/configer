package generation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/configcenter/pkg/define"
	"github.com/configcenter/pkg/util"
	"gopkg.in/yaml.v3"
)

func checkFileContent(infrastructure []byte, rawData map[string][]byte, topicIpRange, topicPortRange, tcpPortRange []string) error {
	//校验ip和port
	if len(topicIpRange)%2 != 0 || len(topicPortRange)%2 != 0 || len(tcpPortRange)%2 != 0 ||
		len(topicIpRange) == 0 || len(topicPortRange) == 0 || len(tcpPortRange) == 0 {
		return errors.New("err topicIpRange or topicPortRange or tcpPortRange length")
	}
	for _, ip := range topicIpRange {
		legal := isIpv4Legal(ip)
		if !legal {
			return errors.New(fmt.Sprintf("err ip input:%s", ip))
		}
	}
	for _, port := range topicPortRange {
		_, legal := checkPort(port)
		if !legal {
			return errors.New(fmt.Sprintf("err port input:%s", port))
		}
	}
	//校验infrastructure的内容
	err := checkInfra(infrastructure)
	if err != nil {
		return err
	}
	for filePath, fileContent := range rawData {
		if strings.Contains(filePath, define.Deployment) {
			//部署信息文件结构问题检查，下同
			pathSli := strings.SplitN(filePath, "/", 7)
			if len(pathSli) != 7 {
				return fmt.Errorf("err deployment file path of:%s, please checkout input path or pkg format", filePath)
			}
			if pathSli[4] != define.DeploymentFlag {
				return fmt.Errorf("err deployment file path of:%s, differ from standard path with flag: %s", filePath, define.DeploymentFlag)
			}
			//部署信息
		} else if strings.Contains(filePath, define.Service) {
			pathSli := strings.SplitN(filePath, "/", 6)
			if len(pathSli) < 6 {
				return fmt.Errorf("err service file path of:%s, please checkout input path or pkg format", filePath)
			}
			if pathSli[4] != define.ServiceFlag {
				return fmt.Errorf("err service file path of:%s, differ from standard path with flag: %s", filePath, define.ServiceFlag)
			}
			//服务信息内容检查
			err = checkService(filePath, fileContent)
			if err != nil {
				return err
			}
		} else if strings.Contains(filePath, define.Template) {
			pathSli := strings.SplitN(filePath, "/", 6)
			if len(pathSli) < 5 || pathSli[4] != define.TemplateFlag {
				return fmt.Errorf("err template file path of:%s, please checkout input path or pkg format", filePath)
			}
		} else if !strings.Contains(filePath, define.Perms) {
			//处理剩余情况
			return fmt.Errorf("a file should not be here with filepath: %s", filePath)
		}
	}

	return nil
}

func checkInfra(infrastructure []byte) error {
	var infraStruct InfraMain
	err := yaml.Unmarshal(infrastructure, &infraStruct)
	if err != nil {
		return fmt.Errorf("decoding infrastructure err: %s", err.Error())
	}
	//校验基础设施信息的键值是否符合预期
	infraDecode, _ := yaml.Marshal(infraStruct)
	keysOk := util.CheckYaml(infrastructure, infraDecode)
	if !keysOk {
		return errors.New("unexpected keys on infrastructure")
	}
	//校验主机名、网卡名、IP是否重复
	hostMap := make(map[string]int)
	for ih, _ := range infraStruct.Host {
		//校验主机名是否重复
		if _, ok := hostMap[infraStruct.Host[ih].HostName]; ok {
			return fmt.Errorf("repeating hostname in infrastructure: %s", infraStruct.Host[ih].HostName)
		}
		//校验绑定的网络名、网卡名、IP是否重复
		netNameMap, adapterMap, IPMap := make(map[string]int), make(map[string]int), make(map[string]int)
		for in, _ := range infraStruct.Host[ih].Network {
			if _, ok := netNameMap[infraStruct.Host[ih].Network[in].Name]; ok {
				return fmt.Errorf("repeating netname in infrastructure on host: %s, net: %s",
					infraStruct.Host[ih].HostName, infraStruct.Host[ih].Network[in].Name)
			}
			if _, ok := adapterMap[infraStruct.Host[ih].Network[in].Adaptor]; ok {
				return fmt.Errorf("repeating adapter in infrastructure on host: %s, adapter: %s",
					infraStruct.Host[ih].HostName, infraStruct.Host[ih].Network[in].Adaptor)
			}
			if _, ok := IPMap[infraStruct.Host[ih].Network[in].Ipv4]; ok {
				return fmt.Errorf("repeating IP in infrastructure on host: %s, IP: %s",
					infraStruct.Host[ih].HostName, infraStruct.Host[ih].Network[in].Ipv4)
			}
			if !isIpv4Legal(infraStruct.Host[ih].Network[in].Ipv4) {
				return fmt.Errorf("illegal ipv4 input:%s, on host:%s", infraStruct.Host[ih].Network[in].Ipv4, infraStruct.Host[ih].HostName)
			}
			netNameMap[infraStruct.Host[ih].Network[in].Name] = 0
			adapterMap[infraStruct.Host[ih].Network[in].Adaptor] = 0
			IPMap[infraStruct.Host[ih].Network[in].Ipv4] = 0
		}
		hostMap[infraStruct.Host[ih].HostName] = 0
	}
	return nil
}

func checkService(path string, data []byte) error {
	var srvStruct SrvMain
	err := yaml.Unmarshal(data, &srvStruct)
	if err != nil {
		return fmt.Errorf("load service.yaml err, file path:%s, err:%s", path, err.Error())
	}
	//校验服务信息的键值是否符合预期
	srvDecode, _ := yaml.Marshal(srvStruct)
	srvOk := util.CheckYaml(data, srvDecode)
	if !srvOk {
		return fmt.Errorf("unexpected keys on %s, please checkout carefully", path)
	}
	subNetMap, pubNetMap := make(map[string]int), make(map[string]int)
	for is, _ := range srvStruct.PubTopic {
		//校验网络名是否重复
		if _, ok := pubNetMap[srvStruct.PubTopic[is].NetName]; ok {
			return fmt.Errorf("repeating net in pub, path: %s, netName: %s", path, srvStruct.PubTopic[is].NetName)
		}
		pubNetMap[srvStruct.PubTopic[is].NetName] = 0
		pubSetTopicMap, pubBizTopicMap := make(map[string]int), make(map[string]int)
		for in, _ := range srvStruct.PubTopic[is].SetTopic {
			if _, ok := pubSetTopicMap[srvStruct.PubTopic[is].SetTopic[in].TpcName]; ok {
				return fmt.Errorf("repeating setTopicName in pub, path: %s, topicName: %s",
					path, srvStruct.PubTopic[is].SetTopic[in].TpcName)
			}
			pubSetTopicMap[srvStruct.PubTopic[is].SetTopic[in].TpcName] = 0
		}
		for in, _ := range srvStruct.PubTopic[is].BizTopic {
			if _, ok := pubBizTopicMap[srvStruct.PubTopic[is].BizTopic[in].TpcName]; ok {
				return fmt.Errorf("repeating bizTopicName in pub, path: %s, topicName: %s",
					path, srvStruct.PubTopic[is].BizTopic[in].TpcName)
			}
			pubBizTopicMap[srvStruct.PubTopic[is].BizTopic[in].TpcName] = 0
		}
	}
	for is, _ := range srvStruct.SubTopic {
		//校验网络名是否重复
		if _, ok := subNetMap[srvStruct.SubTopic[is].NetName]; ok {
			return fmt.Errorf("repeating net in sub, path: %s, netName: %s", path, srvStruct.SubTopic[is].NetName)
		}
		subNetMap[srvStruct.SubTopic[is].NetName] = 0
		subSetTopicMap, subBizTopicMap := make(map[string]int), make(map[string]int)
		for in, _ := range srvStruct.SubTopic[is].SetTopic {
			if _, ok := subSetTopicMap[srvStruct.SubTopic[is].SetTopic[in].TpcName]; ok {
				return fmt.Errorf("repeating setTopicName in sub, path: %s, topicName: %s",
					path, srvStruct.SubTopic[is].SetTopic[in].TpcName)
			}
			subSetTopicMap[srvStruct.SubTopic[is].SetTopic[in].TpcName] = 0
		}
		for in, _ := range srvStruct.SubTopic[is].BizTopic {
			if _, ok := subBizTopicMap[srvStruct.SubTopic[is].BizTopic[in].TpcName]; ok {
				return fmt.Errorf("repeating bizTopicName in sub, path: %s, topicName: %s",
					path, srvStruct.SubTopic[is].BizTopic[in].TpcName)
			}
			subBizTopicMap[srvStruct.SubTopic[is].BizTopic[in].TpcName] = 0
		}
	}
	return nil
}
