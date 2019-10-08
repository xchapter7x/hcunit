# hcunit
Helm Chart Unit: helps to unit test rendering of your templates using policies

## Sample usage
```
$> hcunit mychartdir
Running Policy `values_invalid_test.rego` with fixture data in `values_invalid_test.yml`... [PASS]
Running Policy `values_valid_test.rego` with fixture data in `values_valid_test.yml`... [PASS]
Running Policy `values_scenariob_test.rego` with fixture data in `values_scenariob_test.yml`... [FAIL]
--------------------------------------------------------------------------------

data.authz.test_post_allowed: FAIL (607ns)
data.authz.test_get_anonymous_denied: PASS (288ns)
data.authz.test_get_user_allowed: PASS (346ns)
data.authz.test_get_another_user_denied: PASS (365ns)
--------------------------------------------------------------------------------
PASS: 3/4
FAIL: 1/4
```

## Functionality
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

