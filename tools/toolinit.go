package tools

var headers map[string]string

func init() {
	headers = map[string]string{}

	headers["cookie"] = `_ga=GA1.1.390781912.1665500669; _ga_MF5DDQ4TF9=GS1.1.1665828318.2.0.1665828318.0.0.0; fork_session=f0e57de30d2b0257; mp_f39218c184e602888e26ea53b00435dd_mixpanel={"distinct_id": "yuhongda0315@163.com","$device_id": "183c7f26f01ab7-0c84907d721d63-19525635-13c680-183c7f26f02c6e","$initial_referrer": "$direct","$initial_referring_domain": "$direct","$user_id": "yuhongda0315@163.com","g_team_name": "grtd","g_version": "1.21.0"}; intercom-id-d1key3b8=dfdfa42f-ebfb-4b69-8152-15773e575cae; intercom-session-d1key3b8=; _ga_G812B88X1Y=GS1.1.1665847698.5.0.1665847698.0.0.0`
	headers["content-type"] = "application/json"
	headers["content-length"] = "39"
}
