# CGroup Power saver

Save power by limiting application's CPU usage.


Here is an example of configuration
```yaml

rules:
  # Limit 10% cpu usage for chrome
  - process: ".*[Cc][Hh][Rr][Oo][Mm][Ee].*" # Reges for process name
    resource: "cpu"  # Type of resource to limit. Only cpu is supported now
    amount: 10  # Percentage of CPU that the processes are allowed to use
    when: "battery" # "battery" or "always" This field does not work for now
```