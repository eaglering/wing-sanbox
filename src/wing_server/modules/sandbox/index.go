package sandbox

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"wing_server/modules/sandbox/config"
	pb "wing_server/modules/sandbox/proto"
)

type server struct{}

func (s *server) Compile(ctx context.Context, in *pb.Input) (*pb.Output, error) {
	script := config.Sandbox[in.Language]
	if script == nil {
		return nil, errors.New("尚不支持的编译语言")
	}
	if in.Watchdog == 0 {
		in.Watchdog = 600
	}
	if in.CpuShares == 0 {
		in.CpuShares = 512 // CPU权重默认1024，这里分配一半
	}
	if in.Memory == "" {
		in.Memory = "200m"
	}
	if in.MemorySwap == "" {
		in.MemorySwap = "300m"
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(in.Language + in.Data))
	cipher := hex.EncodeToString(md5Ctx.Sum(nil))
	filename := fmt.Sprintf("/tmp/%v", cipher)
	_, err := os.Stat(filename)
	if err != nil {
		err = ioutil.WriteFile(filename, []byte(in.Data), 0644)
		if err != nil {
			return nil, errors.New("创建文件失败")
		}
		defer os.Remove(filename)
	}
	command := fmt.Sprintf("/usr/local/bin/wing.sh %d %v %v", in.Watchdog, filename, script.Filename)
	switch in.Language {
	case "mysql":
		command = fmt.Sprintf("%v -u mysql", command)
	case "nodejs":
		command = fmt.Sprintf("%v -e NODE_PATH=/usr/local/lib/node_modules", command)
	}

	command = fmt.Sprintf("%v -c %v -m %v --memory-swap %v %v%v \"/usr/local/bin/script.sh %v %v %v %v\"",
		command, in.CpuShares, in.Memory, in.MemorySwap, config.DockerAddress, script.DockerName, script.CompilerName,
		script.Filename, script.OutputCmd, script.Arguments)
	log.Println(command)

	cmd := exec.Command("/bin/bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil, errors.New("执行脚本失败")
	}
	response := string(out)
	if response == "" || !strings.Contains(response, config.EOF) {
		return nil, errors.New("执行超时")
	}
	pos := strings.LastIndex(response, config.EOF)
	data := response[0:pos]
	runtime, err := strconv.ParseFloat(response[pos+len(config.EOF):], 64)
	if err != nil {
		runtime = 0.00
	}
	return &pb.Output{
		Language: script.Language,
		Runtime:  runtime,
		Data:     data,
	}, nil
}

func Register(grpc *grpc.Server) {
	strict := os.Getenv("STRICT")
	if strict == "true" {
		for _, script := range config.Sandbox {
			inspect := fmt.Sprintf("docker inspect %v%v", config.DockerAddress, script.DockerName)
			cmd := exec.Command("/bin/bash", "-c", inspect)
			_, err := cmd.Output()
			if err != nil {
				pull := fmt.Sprintf("docker pull %v%v", config.DockerAddress, script.DockerName)
				cmd = exec.Command("/bin/bash", "-c", pull)
				out, err := cmd.Output()
				if err != nil {
					log.Println(string(out))
				}
			}
		}
	}
	pb.RegisterSandboxServer(grpc, &server{})
}
