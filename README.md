# hcunit
[![CircleCI](https://circleci.com/gh/xchapter7x/hcunit.svg?style=svg)](https://circleci.com/gh/xchapter7x/hcunit)


Helm Chart Unit: helps to unit test rendering of your templates using policies

## Download Binaries
https://github.com/xchapter7x/hcunit/releases/latest

## About hcunit
- Uses [OPA and Rego](https://www.openpolicyagent.org/) to evaluate the yaml to see if it meets your expectations
- By convention hcunit will run any rules in your given rego file or recursively in a given directory as long as that rule takes the form `expect ["..."] { ... } `. it is a good idea to define the hash value within the rule so it prints during a `--verbose` call 
- Your policy rules will have access to a input object. This object will be a hashmap of your rendered templates, with the hash being the filename, and the value being an object representation of the rendered yaml. It will also contain a hash for the NOTES file, which will be a string. 
- uses helm's packages to render the templates so, it should yield identical output as the `helm template` command


## Options
```bash
-> % hcunit --help
Usage:
  hcunit [OPTIONS] <eval | render | version>

Help Options:
  -h, --help  Show this help message

Available commands:
  eval     evaluate a policy on a chart + values
  render   Render a template yaml
  version  display version info
```



## Sample usage
```bash
000@000-000 [00:00:00] [helm-charts/concourse] [master *]
-> % cat policy/testing.rego
───────┬───────────────────────────────────────────────────────────────
       │ File: policy/testing.rego
───────┼───────────────────────────────────────────────────────────────
   1   │ package main
   2   │
   3   │ expect [msg] {
   4   │   msg = "noop pass rule"
   5   │   true
   6   │ }
   7   │
   8   │ expect [msg] {
   9   │   msg = "we should have values and secrets"
  10   │   input["values.yaml"]
  11   │   n = input["web-secrets.yaml"].metadata.name
  12   │   n == "hcunit-name-web"
  13   │ }
───────┴───────────────────────────────────────────────────────────────

000@000-000 [00:00:00] [helm-charts/concourse] [master *]
-> % hcunit eval -t templates/ -c values.yaml -p policy/testing.rego
[PASS] Your policy rules have been run successfully!

000@000-000 [00:00:00] [helm-charts/concourse] [master *]
-> % cat policy/testing_fail.rego
───────┬───────────────────────────────────────────────────────────────
       │ File: policy/testing_fail.rego
───────┼───────────────────────────────────────────────────────────────
   1   │ package main
   2   │
   3   │ expect [msg] {
   4   │   msg = "noop pass rule"
   5   │   true
   6   │ }
   7   │
   8   │ expect [msg] {
   9   │   msg = "we should have values and secrets"
  10   │   input["values.yaml"]
  11   │   n = input["web-secrets.yaml"].metadata.name
  12   │   n == "WRONGNAME"
  13   │ }
───────┴───────────────────────────────────────────────────────────────

000@000-000 [00:00:00] [helm-charts/concourse] [master *]
-> % hcunit eval -t templates/ -c values.yaml -p policy/testing_fail.rego
[FAIL] Your policy rules are violated in your rendered output!
your policy failed

```













## Future Suite Functionality
hcunit will traverse the `test` directory and for each _test.yml & _test.rego pair will use the `.yml` file as a values input for rendering and run the corresponding rego policy against the rendered templates. the input object made available in the rego policy will be a hashmap with keys using the paths of the rendered templates and the values file. The corresponding value in the hashmap being the object representation of the files.


## Chart Dir Convention
```
chart
│   README.md    
│
└───templates
│   │   NOTES.txt
│   │   _helpers.tpl
│   │   web-deployment.yaml
│   │   web-ingress.yaml
│   
└───test
    │   values_invalid_test.yml
    │   values_invalid_test.rego
    │   values_valid_test.yml
    │   values_valid_test.rego 
    │   values_scenariob_test.yml
    │   values_scenariob_test.rego
```

