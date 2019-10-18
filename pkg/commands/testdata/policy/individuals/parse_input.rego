package main

expect ["input should contain valid object representation of rendered template"] {
  input["hcunit/testdata/templates/something.yml"].kind
}
