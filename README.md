# hcunit
[![CircleCI](https://circleci.com/gh/xchapter7x/hcunit.svg?style=svg)](https://circleci.com/gh/xchapter7x/hcunit)


Helm Chart Unit: helps to unit test rendering of your templates using policies

## Usage as a Helm Plugin:

```bash
$> echo "install the latest version of the plugin"
$> helm plugin install https://github.com/xchapter7x/hcunit/releases/latest/download/hcunit_plugin.tgz
Installed plugin: unit

$> echo "you might have have to make the plugin binaries executable"
$> helm env | grep "HELM_PLUGIN" | awk -F"=" '{print $2}' | awk -F\" '{print "chmod +x "$2"/hcunit_plugin/hcunit*"}' |sh

$> echo "lets run some tests of our templates' logic"
$> helm unit -t templates -c policy/values_toggle_on.yaml -p policy/testing_toggle_on.rego
[PASS] Your policy rules have been run successfully!

$> echo "lets explore the available flags for the plugin call"
$> helm unit --help
Usage:
  hcunit_osx [OPTIONS] eval [eval-OPTIONS]

given a OPA/Rego Policy one can evaluate if the rendered templates of a chart using a given values file meet the defined rules of the policy or not

Help Options:
  -h, --help           Show this help message

[eval command options]
      -t, --template=  path to yaml template you would like to render
      -c, --values=    path to values file you would like to use for rendering
      -p, --policy=    path to rego policies to evaluate against rendered templates
      -n, --namespace= policy namespace to query for rules
      -v, --verbose    prints tracing output to stdout
      
```



## Usage as a Standalone CLI... Download Binaries
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













