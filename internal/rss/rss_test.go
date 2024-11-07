package rss

import "testing"

func Test_getRSS(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "RSS1",
			args: args{
				url: "https://habr.com/ru/rss/hub/go/all/?fl=ru",
			},
		},
		{
			name: "RSS2",
			args: args{
				url: "https://habr.com/ru/rss/best/daily/?fl=ru",
			},
		},
		{
			name: "RSS3",
			args: args{
				url: "https://cprss.s3.amazonaws.com/golangweekly.com.xml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := getRSS(tt.args.url)
			if err != nil {
				t.Error(err)
			}
			t.Log(list[0].Title)
			t.Log(list[1].Title)
			t.Log(list[2].Title)
		})
	}
}
