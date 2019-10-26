package main

allow ["input object should provide values in hash"] {
  input["values.yml"]
}
deny ["input object should provide values in hash"] {
  input["values.yml"]
}
