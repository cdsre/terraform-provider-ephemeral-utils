ephemeral "vault_kv_secret_v2" "foo" {
  mount = "secret"
  name  = "foobar"
}

resource "ephemeral-utils_revealer" "foo" {
  data_wo = ephemeral.vault_kv_secret_v2.foo.custom_metadata["desc"]
}
