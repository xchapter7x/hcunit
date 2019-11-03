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



## Notes on Syntax and Rego

Rego is a Policy Language for the Open Policy Agent eco system. We use rego here as our testing DSL. Any rego rule which is an `assert` or `expect` will get executed and must evaluated to true. The gist is that everything between the `{}` is a `rule`. Everything between `{}` should evaluate to `true`. Assignments yield true, and if any statement in the `{}` block is `false` then the entire rule will return `false` and therfore fail our test case.

For more information you can try: https://www.openpolicyagent.org/docs/latest/#rego

Or 

for a Online playground: https://play.openpolicyagent.org/





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
   3   │ assert [behavior] {
   4   │   behavior = "this should always be true b/c its true"
   5   │   true
   6   │ }
   7   │
   8   │ assert [behavior] {
   9   |     behavior = "when web is enabled then namespace is toggled on"
   10  |     "true" == input["values.yaml"].web.enabled
   11  │     "Namespace" == input["namespace.yaml"].kind
   12  │ }
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
   3   │ assert [behavior] {
   4   │   behavior = "this should always be true b/c its true"
   5   │   false
   6   │ }
   7   │
   8   │ assert [behavior] {
   9   |     behavior = "when web is enabled then namespace is toggled on"
   10  |     "true" == input["values.yaml"].web.enabled
   11  │     "NamespaceWrongKind" == input["namespace.yaml"].kind
   12  │ }
───────┴───────────────────────────────────────────────────────────────

000@000-000 [00:00:00] [helm-charts/concourse] [master *]
-> % hcunit eval -t templates/ -c values.yaml -p policy/testing_fail.rego
[FAIL] Your policy rules are violated in your rendered output!
your policy failed

```


## About hcunit
- Uses [OPA and Rego](https://www.openpolicyagent.org/) to evaluate the yaml to see if it meets your expectations
- By convention hcunit will run any rules in your given rego file or recursively in a given directory as long as that rule takes the form `assert ["some behavior"] { ... } ` or `expect ["some other behavior"] { ... } `. it is a good idea to define the hash value within the rule so it prints during a `--verbose` call 
- Your policy rules will have access to a input object. This object will be a hashmap of your rendered templates, with the hash being the filename, and the value being an object representation of the rendered yaml. It will also contain a hash for the NOTES file, which will be a string. 
- uses helm's packages to render the templates so, it should yield identical output as the `helm template` command






