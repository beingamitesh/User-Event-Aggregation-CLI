package main

import (
	"reflect"
	"testing"
)

func Test_checkIfNewEventHashExists(t *testing.T) {
	type args struct {
		eventHashes map[string]interface{}
		event       UserEvent
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Testing if event hash already exists in the map",
			args: args{
				eventHashes: map[string]interface{}{
					"0ae9db6838cff2fa3f64c723f3bb1b6c200d1d8bc6fd96cd6abbf31dab731a2f": nil,
					"1b179b7b75b84d2cc5277a424dba4719f62e1653c0410b7902130590ceb2bd8e": nil,
					"41cd11117720f96ebcc7c44b3cf7877dbd6160813607901e0ecabd267778adb0": nil,
					"5cb854dd619e5d3b133f818b2f85e3e68271006c3f57c470dc998b06b73588e5": nil,
					"807d00814f33ccb875b85cf8067c48d21902ca73ec5954e6f199eaeb2963d721": nil,
					"f6e774040f4be937a766d51b2638b128e11386870c89ef8bb0ad89bb381e743f": nil,
				},
				event: UserEvent{
					UserID:    1,
					EventType: "post",
					Timestamp: 1672444800,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Testing if event hash does not exists in the map",
			args: args{
				eventHashes: map[string]interface{}{
					"0ae9db6838cff2fa3f64c723f3bb1b6c200d1d8bc6fd96cd6abbf31dab731a2f": nil,
					"1b179b7b75b84d2cc5277a424dba4719f62e1653c0410b7902130590ceb2bd8e": nil,
					"41cd11117720f96ebcc7c44b3cf7877dbd6160813607901e0ecabd267778adb0": nil,
					"5cb854dd619e5d3b133f818b2f85e3e68271006c3f57c470dc998b06b73588e5": nil,
					"807d00814f33ccb875b85cf8067c48d21902ca73ec5954e6f199eaeb2963d721": nil,
					"f6e774040f4be937a766d51b2638b128e11386870c89ef8bb0ad89bb381e743f": nil,
				},
				event: UserEvent{
					UserID:    3,
					EventType: "post",
					Timestamp: 1672531203,
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfNewEventHashExists(tt.args.eventHashes, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfNewEventHashExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfNewEventHashExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateSummary(t *testing.T) {
	type args struct {
		summaries                     *DailySummaries
		event                         UserEvent
		UserIdTimestampToEventTypeMap map[int]map[string]map[string]int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Testing for the first event entry in summary",
			args: args{
				summaries: &DailySummaries{},
				event: UserEvent{
					UserID:    1,
					EventType: "post",
					Timestamp: 1672444800,
				},
				UserIdTimestampToEventTypeMap: make(map[int]map[string]map[string]int),
			},
		},
		{
			name: "Testing for the second event entry for same user, same timestamp(date) but different event type",
			args: args{
				summaries: &DailySummaries{
					DailySummary{
						"date":   "2022-12-31",
						"post":   1,
						"userId": 1,
					},
				},
				event: UserEvent{
					UserID:    1,
					EventType: "likeReceived",
					Timestamp: 1672444801,
				},
				UserIdTimestampToEventTypeMap: map[int]map[string]map[string]int{
					1: map[string]map[string]int{
						"2022-12-31": map[string]int{
							"post": 1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateSummary(tt.args.summaries, tt.args.event, tt.args.UserIdTimestampToEventTypeMap)
		})
	}
}

func Test_getOrCreateSummary(t *testing.T) {
	type args struct {
		summaries *DailySummaries
		userID    int
		date      string
	}
	tests := []struct {
		name string
		args args
		want DailySummary
	}{
		{
			name: "Testing if no summary entry exists",
			args: args{
				summaries: &DailySummaries{},
				userID:    1,
				date:      "2022-12-31",
			},
			want: DailySummary{
				"date":   "2022-12-31",
				"userId": 1,
			},
		},
		{
			name: "Testing if summary entry exists",
			args: args{
				summaries: &DailySummaries{
					DailySummary{
						"date":   "2022-12-31",
						"post":   1,
						"userId": 1,
					},
				},
				userID: 1,
				date:   "2022-12-31",
			},
			want: DailySummary{
				"date":   "2022-12-31",
				"post":   1,
				"userId": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOrCreateSummary(tt.args.summaries, tt.args.userID, tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOrCreateSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateSummaryFields(t *testing.T) {
	type args struct {
		summary                       DailySummary
		event                         UserEvent
		UserIdTimestampToEventTypeMap map[int]map[string]map[string]int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Testing for the first entry in summary",
			args: args{
				summary: nil,
				event: UserEvent{
					UserID:    1,
					EventType: "post",
					Timestamp: 1672444800,
				},
				UserIdTimestampToEventTypeMap: make(map[int]map[string]map[string]int),
			},
		},
		{
			name: "Testing for the second entry in summary",
			args: args{
				summary: DailySummary{
					"date":   "2022-12-31",
					"post":   1,
					"userId": 1,
				},
				event: UserEvent{
					UserID:    1,
					EventType: "likeReceived",
					Timestamp: 1672444801,
				},
				UserIdTimestampToEventTypeMap: map[int]map[string]map[string]int{
					1: map[string]map[string]int{
						"2022-12-31": map[string]int{
							"post": 1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateSummaryFields(tt.args.summary, tt.args.event, tt.args.UserIdTimestampToEventTypeMap)
		})
	}
}

func Test_hashFunc(t *testing.T) {
	type args struct {
		event   UserEvent
		hashMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Testing successful hashing",
			args: args{
				event: UserEvent{
					UserID:    1,
					EventType: "post",
					Timestamp: 1672444800,
				},
				hashMap: make(map[string]interface{}),
			},
		},
		{
			name: "Testing successful hashing - 2",
			args: args{
				event: UserEvent{
					UserID:    1,
					EventType: "likeReceived",
					Timestamp: 1672444801,
				},
				hashMap: map[string]interface{}{
					"0ae9db6838cff2fa3f64c723f3bb1b6c200d1d8bc6fd96cd6abbf31dab731a2f": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := hashFunc(tt.args.event, tt.args.hashMap); (err != nil) != tt.wantErr {
				t.Errorf("hashFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
