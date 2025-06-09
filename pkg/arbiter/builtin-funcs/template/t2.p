timerange_prv_10m = dql_timerange_get()
dur = 600000 # (ms) 10 * 60 * 1 * 1000ms

timerange_prv_10m[0] -= dur
timerange_prv_10m[1] -= dur

printf("prev 10min: %v, timerange: %v\n" ,timerange_prv_10m, dql_timerange_get())

pods_alert = []

## ------------- check cpu usage -------------

## query history object data
#
v = dql("OH::`kubelet_pod`:(pod_name, namespace) {cpu_usage_base_limit > 10} by pod_name, namespace", time_range=timerange_prv_10m)

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


## query object
#
v = dql("O::`kubelet_pod`:(pod_name, namespace, cpu_usage_base_limit) {cpu_usage_base_limit > 10} by pod_name, namespace")

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
            pods_alert = append(pods_alert, {
                "pod_name": pod,
                "namespace": ns,
                "cpu_usage": cpu,
            })
        }
    }
}


printf("%v\n", pods_alert)