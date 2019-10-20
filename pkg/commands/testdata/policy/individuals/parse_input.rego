package main

expect ["input should contain valid object representation of rendered template"] {
  k := input["something.yml"].kind
  k == "Ingress"
}
