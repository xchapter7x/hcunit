package main

expect ["another force failure 123"] {
  not true
}

expect ["force failure 456"] {
  not true
}

expect ["some things pass 789"] {
  true
}

expect ["another passing case 10 11 12"] {
  true
}
