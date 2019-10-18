package main

expect ["another force failure"] {
  not true
}

expect ["force failure"] {
  not true
}

expect ["some things pass"] {
  true
}

expect ["another passing case"] {
  true
}
