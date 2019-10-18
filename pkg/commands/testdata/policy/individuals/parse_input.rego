package main

expect ["input should contain valid object representation of rendered template"] {
  k := input["hcunit/testdata/templates/something.yml"].kind
  k == "Ingress"
}
