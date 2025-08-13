terraform {
  required_version = ">= 1.0"
}

# Simple local file resource for hello world
resource "local_file" "hello" {
  filename = "${path.module}/hello.txt"
  content  = "Hello, World from OpenTofu!"
}

# Output the file content
output "hello_content" {
  value = local_file.hello.content
}

output "file_path" {
  value = local_file.hello.filename
}
