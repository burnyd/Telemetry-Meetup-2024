#### Capabilities of a device
gnmic -a 172.20.20.2:6030 -u admin -p admin capabilities --insecure

#### Get of all paths
gnmic -a 172.20.20.2:6030 -u admin -p admin get --path "/" --insecure

#### Streaming interface stats
gnmic -a 172.20.20.2:6030 -u admin -p admin subscribe --path "openconfig:/interfaces/interface/state/counters" --insecure

#### AFTs
gnmic -a 172.20.20.2:6030 -u admin -p admin subscribe --path "network-instances/network-instance/afts/" --insecure

#### eos_native paths
gnmic -a 172.20.20.2:6030 -u admin -p admin --insecure get --path 'eos_native:/Kernel/proc/cpu/utilization/total' --insecure

#### Cli Origin
gnmic -a 172.20.20.2:6030 -u admin -p admin --insecure get --path 'cli:/show version' --insecure