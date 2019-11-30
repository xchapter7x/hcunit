package main

assert ["values object should have multiple file values"] {
  true == input["values"]["uiIngress"]["enabled"]
}
