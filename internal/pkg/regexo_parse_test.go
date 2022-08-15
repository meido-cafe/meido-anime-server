package pkg

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetEpisode(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name        string
		args        args
		wantEpisode int64
		wantErr     bool
	}{
		{"", args{title: "[桜都字幕组] 不死者之王 第四季 / OVERLORD Ⅳ [05][1080p][简体内嵌]"}, 5, false},
		{"", args{title: "[桜都字幕组] 不死者之王 第季 / OVERLORD Ⅳ [ 05 ][1080p][简体内嵌]"}, 5, false},
		{"", args{title: "[桜都字幕组] 不死者之王 第十三季 / OVERLORD Ⅳ [ 5 ][1080p][简体内嵌]"}, 5, false},
		{"", args{title: "[桜都字幕组] 不死者之王 第三四二季 / OVERLORD Ⅳ [5][1080p][简体内嵌]"}, 5, false},
		{"", args{title: "[桜都字幕组] 不死者之王 第三十四季 / OVERLORD Ⅳ [15][1080p][简体内嵌]"}, 15, false},
		{"", args{title: "[桜都字幕组] 不死者之王 第四期 / OVERLORD Ⅳ [ 35 ][1080p][简体内嵌]"}, 35, false},
		{"", args{title: "[桜都字幕组] 不死者之王 s4 / OVERLORD Ⅳ 15 [1080p][简体内嵌]"}, 15, false},
		{"", args{title: "[桜都字幕组] 不死者之王 s12 / OVERLORD Ⅳ [1080p][简体内嵌]"}, 0, true},
		{"", args{title: "[桜都字幕组] 不死者之王 S4 / OVERLORD Ⅳ 06 [1080p][简体内嵌]"}, 6, false},
		{"", args{title: "[桜都字幕组] 不死者之王 s04 / OVERLORD Ⅳ 8 [1080p][简体内嵌]"}, 8, false},
		{"第x集", args{title: "[桜都字幕组] 不死者之王 s04 / OVERLORD Ⅳ 第08集 [1080p][简体内嵌]"}, 8, false},
		{"第x话 #2", args{title: "[桜都字幕组] 不死者之王 s04 / OVERLORD Ⅳ 第08话 [1080p][简体内嵌]"}, 8, false},
		{"中文数字1", args{title: "[桜都字幕组] 不死者之王 s04 / OVERLORD Ⅳ 第八话 [1080p][简体内嵌]"}, 8, false},
		{"中文数字2", args{title: "[桜都字幕组] 不死者之王 s04 / OVERLORD Ⅳ 第八集 [1080p][简体内嵌]"}, 8, false},
		{"中文数字 error", args{title: "[桜都字幕组] 不死者之王 s04 / OVERLORD Ⅳ 8 [1080p][简体内嵌]"}, 8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEpisode, err := GetEpisode(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEpisode != tt.wantEpisode {
				t.Errorf("GetEpisode() gotEpisode = %v, want %v", gotEpisode, tt.wantEpisode)
			}
		})
	}
}

func TestGetSeason(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name       string
		args       args
		wantSeason int64
		wantStatus int
		wantErr    bool
	}{
		{"<10", args{title: "OVERLORD 第四季"}, 4, 0, false},
		{"无空格", args{title: "OVERLORD第四季"}, 1, 1, false},
		{"无数字", args{title: "OVERLORD 第季"}, 1, 1, false},
		{"10+", args{title: "OVERLORD 第十三季"}, 13, 0, false},
		{">20", args{title: "OVERLORD 第三十四季"}, 34, 0, false},
		{"季信息格式错误", args{title: "OVERLORD 第三四二季"}, 0, 0, true},
		{"期", args{title: "OVERLORD 第四期"}, 4, 0, false},
		{"<10 没有0", args{title: "OVERLORD s4"}, 4, 0, false},
		{">10", args{title: "OVERLORD s12"}, 12, 0, false},
		{">100", args{title: "OVERLORD s12232"}, 12232, 0, false},
		{"罗马数字", args{title: "OVERLORD IV"}, 4, 0, false},
		{"<10 有0", args{title: "OVERLORD s04"}, 4, 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeason, gotStatus, err := GetSeason(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSeason() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSeason != tt.wantSeason {
				t.Errorf("GetSeason() gotSeason = %v, want %v", gotSeason, tt.wantSeason)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetSeason() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func Test_preHandlerTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1",
			args{title: "【喵萌奶茶屋&LoliHouse】 杜鹃的婚约 / Kakkou no Iinazuke]【06】[WebRip 1080p HEVC-10bit AAC][简繁内封字幕]"},
			"[喵萌奶茶屋&LoliHouse] 杜鹃的婚约 / Kakkou no Iinazuke][06][WebRip 1080p HEVC-10bit AAC][简繁内封字幕]",
		},
		{
			"2",
			args{title: "[   喵萌奶茶屋&LoliHouse ] 杜鹃的婚约 / Kakkou no Iinazuke     ][06][   WebRip 1080p HEVC-10bit AAC]   [   简繁内封字幕]"},
			"[喵萌奶茶屋&LoliHouse] 杜鹃的婚约 / Kakkou no Iinazuke][06][WebRip 1080p HEVC-10bit AAC]   [简繁内封字幕]",
		},
		{
			"3",
			args{title: "[喵萌奶茶屋&LoliHouse] 杜鹃的婚约 / Kakkou no Iinazuke][ 06 ][WebRip 1080p HEVC-10bit AAC][简繁内封字幕]"},
			"[喵萌奶茶屋&LoliHouse] 杜鹃的婚约 / Kakkou no Iinazuke][06][WebRip 1080p HEVC-10bit AAC][简繁内封字幕]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreHandleTitle(tt.args.title); got != tt.want {
				t.Errorf("preHandlerTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASADAD(t *testing.T) {
	title := "[桜都字幕组] 不死者之王 S4 / OVERLORD Ⅳ [05][1080p][简体内嵌]"
	fmt.Println(title)
	title = strings.ToLower(title)
	fmt.Println(title)
}
