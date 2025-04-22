# v, ok = dql("T::sre(`.*`):(host, service, span_id, status) limit 2")
# if ok {
#     printf("dql query status code: %d\n", v["status_code"])
#     if ok && v["status_code"] == 200 {
#         for pt in v["points"] {
#             trigger(pt)
#         }
#     }
# } else {
#     printf("dql query execution failed\n")
# }

v, ok = dql("M::cpu:(usage_user as user, usage_total as total) limit 2 slimit 2")
#v, ok = dql("T::re(`.*`):(service, span_id, status) by host limit 3")
#v, ok = dql("M::cpu limit 2 slimit 2 by host")
printf("%v", v)
