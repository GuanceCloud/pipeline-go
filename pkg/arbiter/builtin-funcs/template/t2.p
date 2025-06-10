timerange_prv_10m = dql_timerange_get()
dur = 600000 # (ms) 10 * 60 * 1 * 1000ms

timerange_prv_10m[0] -= dur
timerange_prv_10m[1] -= dur

printf("prev 10min: %v, timerange: %v\n" ,timerange_prv_10m, dql_timerange_get())

## ------------- check cpu usage -------------

## query history pod cpu usage
#
v = dql("OH::`kubelet_pod`:(pod_name, namespace) {cpu_usage_base_limit > 90} by pod_name, namespace", time_range=timerange_prv_10m)

series_pod_name = dql_series_get(v, "pod_name")
series_namespace = dql_series_get(v, "namespace")
pods_with_ns = {}


## There are multiple groups
#
for i = 0; i < len(series_pod_name); i+=1 {
    pods = series_pod_name[i]
    nss = series_namespace[i]
    
    ## Do not use `i` a loop variable, since the upper loop has been defined
    #
    for j = 0; j < len(nss); j+=1 {
        ns = nss[j]
        if ns in pods_with_ns {
            pods_with_ns[ns] = append(pods_with_ns[ns], pods[j])
        }  else {
            pods_with_ns[ns] = [pods[j]]
        }
    }
}


## query pod cpu usage
#
v = dql("O::`kubelet_pod`:(pod_name, namespace, cpu_usage_base_limit) {cpu_usage_base_limit > 90} by pod_name, namespace")

series_pod_name = dql_series_get(v, "pod_name")
series_namespace = dql_series_get(v, "namespace")
series_usage = dql_series_get(v, "cpu_usage_base_limit")


for i = 0; i < len(series_pod_name); i+=1 {
    pods = series_pod_name[i]
    nss = series_namespace[i]
    usage = series_usage[i]
    
    for j = 0; j < len(nss); j+=1 {
        ns = nss[j]
        pod = pods[j]
        cpu = usage[j] 
        
        if ns in pods_with_ns && pod in pods_with_ns[ns]{
            trigger(
                result = {
                    "cpu_usage": cpu,
                    "pod_name": pod,
                    "namespace": ns,
                    "kind": "pod_cpu_check"
                },
                status = "high",
                dimension_tags = {
                    "pod_name": pod,
                    "namespace": ns
                }
            )
        }
    }
}



## query k8s node memory usage
#
v = dql("OH::`HOST`:(distinct(host) as host) {mem_used_percent > 80}")
tmp_host = dql_series_get(v, "host")

host_li = []
if len(tmp_host) == 1 {
    host_li = tmp_host[0]
}

v = dql("O::`HOST`:(host, mem_used_percent) {mem_used_percent > 80}")
tmp_host = dql_series_get(v, "host")
tmp_mem = dql_series_get(v, "mem_used_percent")

if len(tmp_host) == 1 {
    hosts = tmp_host[0]
    mem = tmp_mem[0]
    for i = 0; i < len(hosts); i += 1 {
        host_name = hosts[i]
        mem_usage = mem[i]
        if host_name in host_li {
            trigger(
                result = {
                    "mem_usage": mem_usage,
                    "host": host_name,
                    "kind": "node_mem_check"
                },
                status = "high",
                dimension_tags = {
                    "host": host_name
                }
            )
        }
    }
}
