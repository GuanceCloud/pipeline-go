url_path = "/api/v1/auth/signin"

time_range = dql_timerange_get()

v = dql(strfmt(
        "R::resource:(distinct(ip) as ip) {resource_url_path = `%s` and resource_status != 200}",
        url_path
    ),
    time_range = time_range)

series_ip = dql_series_get(v, "ip")

printf("Distinct ip: %v; Query time range: %v, %d min\n", 
    series_ip, 
    time_range, 
    (time_range[1] - time_range[0])/1000/60 )

IPs = []
if len(series_ip) == 1 {
    IPs = series_ip[0]
}


risk_ip = {}

for ip in IPs {
    v = dql(strfmt(
            "R::resource:(count(`*`) as count) {resource_url_path = `%s` and ip = `%s` and resource_status != 200}",
            url_path,
            ip
        )
    )
    count = 0
    series_count = dql_series_get(v, "count")
    if len(series_count) == 1 && len(series_count[0]) == 1 {
        count = series_count[0][0]
    }
    if count >= 10 {
        printf("**High risk IP %s, access `%s` %.0f times\n", ip, url_path, count)
        risk_ip[ip] = count
    }
}

if len(risk_ip) > 0 {
    printf("Risk IPs: %v\n", risk_ip)
    trigger(result=risk_ip, dimension_tags={"url_path": url_path}, status="critical")
} else {
    printf("Risk IP not found")
    trigger(result=risk_ip, dimension_tags={"url_path": url_path}, status="info")
}

