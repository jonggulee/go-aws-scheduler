# AWS Lambda Scheduler

AutoScaling, EC2, RDS 서비스를 자동으로 재시작하는 lambda functions 입니다.

## Getting Started
```
func main() {
	lambda.Start(handler)

    // 로컬에서 테스트 진행 시
	handler()
}
```

## Building your function

AWS Lambda에 배포하려면 Linux로 컴파일 후 .zip 파일에 저장되어야 한다.

### For developers on Linux and macOS

```
GOOS=linux GOARCH=amd64 go build -o main main.go && zip main.zip main
```

## Deploying your functions
```
aws lambda update-function-code --function-name {function name} --zip-file fileb://{zip file name}.zip
```
