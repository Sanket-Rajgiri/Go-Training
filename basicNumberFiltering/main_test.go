package main

import (
	"reflect"
	"testing"
)

func Test_returnPrime(t *testing.T) {
	type args struct {
		in0 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Story3",
			args: args{
				in0: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []int{2, 3, 5, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnPrime(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnPrime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_returnEven(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Story1",
			args: args{
				nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []int{2, 4, 6, 8, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnEven(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnEven() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_returnOdd(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Story2",
			args: args{
				nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []int{1, 3, 5, 7, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnOdd(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnOdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_story4(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "Story4",
			args: args{
				nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []int{3, 5, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := story4(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnOddPrime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_story5(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "Story5",
			args: args{
				nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
			want: []int{10, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := story5(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("returnEvenMultiplesof5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_story6(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "Story6",
			args: args{
				nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
			want: []int{15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := story6(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Story6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_story7(t *testing.T) {
	type args struct {
		nums       []int
		conditions []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "Story7.1",
			args: args{
				nums:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				conditions: []string{"odd", "greater than 5", "multiple of 3"},
			},
			want: []int{9, 15},
		},
		{
			name: "Story7.2",
			args: args{
				nums:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				conditions: []string{"less than 6", "multiple of 3"},
			},
			want: []int{6, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := story7(tt.args.nums, tt.args.conditions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("story7() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_story8(t *testing.T) {
	type args struct {
		nums       []int
		conditions []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "Story8.1",
			args: args{
				nums:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				conditions: []string{"prime", "greater than 15", "multiple of 5"},
			},
			want: []int{2, 3, 5, 7, 10, 11, 13, 15, 16, 17, 18, 19, 20},
		},
		{
			name: "Story8.2",
			args: args{
				nums:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				conditions: []string{"less than 6", "multiple of 3"},
			},
			want: []int{1, 2, 3, 4, 5, 6, 9, 12, 15, 18},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := story8(tt.args.nums, tt.args.conditions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("story8() = %v, want %v", got, tt.want)
			}
		})
	}
}
