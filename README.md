# Go_AWS

Golang에서 AWS서비스를 어떻게 활용하는지에 대한 repo

# 살펴볼 AWS 서비스

1. S3 - 블로그, 코드 완
2. SQS - 블로그, 코드 완
3. API Gateway
- API Gateway는 사싱상 코드에서 할 것이 없기 떄문에 따로 블로그로 정리 예정
4. DynamoDB

# app.env 형태
- 특정 파일이 아니라 os.GetEnv가 더 안전

aws_region= <value> -> 사용하는 aws서비스의 region

access_key= <value>

secret_key= <value>

# 레파지토리 설명

해당 레파지토리에서는 go에서 AWS SDK를 사용하여 어떻게 서비스에 접근하는지에 대한 내용입니다.
일부 주석들이 보인다면, 해당 주석들은 후에 블로그를 통해서 정리되고 삭제 될 예정입니다.

일부 기능은 사용을 할 줄 알지만, 좀 더 다양한 옵션들과 다양한 설정들을 알아보고 테스트 하기 위해 만들어진 레파지토리임을 참고바랍니다.
