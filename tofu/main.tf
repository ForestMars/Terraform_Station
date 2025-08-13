terraform {
  required_version = ">= 1.0"
}

resource "local_file" "hello" {
  filename = "${path.module}/hello.txt"
  content = "Hello, World from OpenTofu!"
}

output "hello_content" {
  value = local_file.hello.content
}
