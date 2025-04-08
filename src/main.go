package main

import (
	"fmt"
	"io"
	"log"
	"memberlist-service-discovery/src/discovery"
	"net"
	"os"
)

func getK8sPodAddresses(serviceName string) ([]string, error) {
	// 查询 Kubernetes Headless Service 的 DNS
	ips, err := net.LookupIP(serviceName)
	if err != nil {
		return nil, err
	}

	addresses := []string{}
	for _, ip := range ips {
		addresses = append(addresses, fmt.Sprintf("%s:6789", ip.String()))
	}
	return addresses, nil
}

func main() {
	// 初始化 Memberlist
	nodeName, _ := os.Hostname()
	memberlistConfig := discovery.DefaultConfig(nodeName)
	memberlistConfig.BindAddr = os.Getenv("POD_IP")
	memberlistConfig.BindPort = 6789
	memberlistConfig.LogOutput = io.Discard

	service, err := discovery.NewMemberlistService(memberlistConfig)
	if err != nil {
		log.Fatalf("初始化 Memberlist 失败: %v", err)
	}
	defer service.Shutdown()

	// 动态获取 Kubernetes Pod 地址
	existingNodes, err := getK8sPodAddresses(fmt.Sprintf("%s.svc.cluster.local", os.Getenv("SERVICE_NAME")))
	if err != nil {
		//log.Printf("无法获取现有服务地址: %v", err)
	} else if len(existingNodes) > 0 {
		err = service.Join(existingNodes)
		if err != nil {
			//log.Printf("服务注册失败: %v", err)
		}
	}

	// 定时打印当前成员
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		fmt.Println("已发现服务列表:")
	//		for nodeName, member := range service.Members() {
	//			fmt.Printf("服务名：%s, 服务地址：%s\n", nodeName, member)
	//		}
	//	}
	//}()

	// 阻塞主线程
	select {}
}
