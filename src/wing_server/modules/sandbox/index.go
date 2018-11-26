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
	}
	command := fmt.Sprintf("/usr/local/bin/wing.sh %d %v %v", in.Watchdog, filename, script.Filename)
	switch in.Language {
	case "mysql":
		command = fmt.Sprintf("%v -u mysql", command)
	case "nodejs":
		command = fmt.Sprintf("%v -e NODE_PATH=/usr/local/lib/node_modules", command)
	}

	command = fmt.Sprintf("%v %v%v \"/usr/local/bin/script.sh %v %v %v %v\"",
		command, config.DockerAddress, script.DockerName, script.CompilerName,
		script.Filename, script.OutputCmd, script.Arguments)
	log.Println(command)

	cmd := exec.Command("/bin/bash", "-c", command)
	out, err := cmd.Output()
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
	env := os.Getenv("ENV")
	if env == "production" {
		for _, script := range config.Sandbox {
			inspect := fmt.Sprintf("docker inspect %v%v >/dev/null", config.DockerAddress, script.DockerName)
			cmd := exec.Command("/bin/bash", "-c", inspect)
			out, err := cmd.Output()
			log.Println(err, string(out))
			if err != nil || len(out) > 0 {
				pull := fmt.Sprintf("docker pull %v%v >/dev/null", config.DockerAddress, script.DockerName)
				cmd = exec.Command("/bin/bash", "-c", pull)
				out, err := cmd.Output()
				log.Println(err, string(out))
				if err != nil || len(out) > 0 {
					log.Fatal(err, string(out))
				}
			}
		}
	}
	pb.RegisterSandboxServer(grpc, &server{})
}