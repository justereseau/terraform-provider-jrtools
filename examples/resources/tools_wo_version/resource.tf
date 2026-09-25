ephemeral "random_password" "db" {
  length = 32
}

resource "tools_wo_version" "db_password" {
  value_wo = ephemeral.random_password.db.result
}

resource "aws_db_instance" "example" {
  # ...
  storage_encrypted   = true
  password_wo         = ephemeral.random_password.db.result
  password_wo_version = tools_wo_version.db_password.version
}
