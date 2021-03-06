package config

import "os"

const (
	EOF = "*-SANDBOX::ENDOFOUTPUT-*"
)

type Script struct {
	Language string
	CompilerName string
	Filename string
	OutputCmd string
	Arguments string
	DockerName string
}

var (
	Sandbox = map[string]*Script{
		"python": {
			Language:"Python",
			CompilerName:"python",
			Filename:"file.py",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-python:v2.0",
		},
		"ruby": {
			Language:"Ruby",
			CompilerName:"ruby",
			Filename:"file.rb",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-ruby:v2.0",
		},
		"clojure": {
			Language:"Clojure",
			CompilerName:"clojure",
			Filename:"file.py",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-clojure:v2.0",
		},
		"php": {
			Language:"PHP",
			CompilerName:"php",
			Filename:"file.php",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-php:v2.0",
		},
		"nodejs": {
			Language:"NodeJS",
			CompilerName:"nodejs",
			Filename:"file.js",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-nodejs:v2.0",
		},
		"scala": {
			Language:"Scala",
			CompilerName:"scala",
			Filename:"file.scala",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-scala:v2.0",
		},
		"golang": {
			Language:"Golang",
			CompilerName:"'go run'",
			Filename:"file.go",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-golang:v2.0",
		},
		"cc": {
			Language:"C/C++",
			CompilerName:"'g++ -o /data/a.out'",
			Filename:"file.cpp",
			OutputCmd:"/data/a.out",
			Arguments:"",
			DockerName:"wing-sandbox-cc:v2.0",
		},
		"java": {
			Language:"Java",
			CompilerName:"javac",
			Filename:"file.java",
			OutputCmd:"/usr/local/bin/javaRunner.sh",
			Arguments:"",
			DockerName:"wing-sandbox-java:v2.0",
		},
		"vbnet": {
			Language:"VB.Net",
			CompilerName:"'vbnc -nologo -quiet'",
			Filename:"file.vb",
			OutputCmd:"'mono /data/file.exe'",
			Arguments:"",
			DockerName:"wing-sandbox-vbnet:v2.0",
		},
		"csharp": {
			Language:"C#",
			CompilerName:"gmcs",
			Filename:"file.cs",
			OutputCmd:"'mono /data/file.exe'",
			Arguments:"",
			DockerName:"wing-sandbox-csharp:v2.0",
		},
		"bash": {
			Language:"Bash",
			CompilerName:"/bin/bash",
			Filename:"file.sh",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-bash:v2.0",
		},
		"oc": {
			Language:"Objective-C",
			CompilerName:"gcc",
			Filename:"file.m",
			OutputCmd:" /data/a.out",
			Arguments:"'-o /data/a.out -I/usr/include -L/usr/lib -lobjc -lgnustep-base -Wall -fconstant-string-class=NSConstantString'",
			DockerName:"wing-sandbox-oc:v2.0",
		},
		"mysql": {
			Language:"MySQL",
			CompilerName:"/usr/local/bin/sql_runner.sh",
			Filename:"file.sql",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-mysql:v2.0",
		},
		"perl": {
			Language:"Perl",
			CompilerName:"perl",
			Filename:"file.pl",
			OutputCmd:"",
			Arguments:"",
			DockerName:"wing-sandbox-perl:v2.0",
		},
		"rust": {
			Language:"Rust",
			CompilerName:"'env HOME=/opt/rust /opt/rust/.cargo/bin/rustc'",
			Filename:"file.rs",
			OutputCmd:"/data/a.out",
			Arguments:"'-o /data/a.out'",
			DockerName:"wing-sandbox-rust:v2.0",
		},
	}
	DockerAddress = os.Getenv("DOCKER_ADDRESS")
)