package main

expect ["force passing"] {
  true; trace(sprintf("[TRACE] %s", [input]))
}

expect ["another passing case"] {
  true
}
