package util

import "testing"

func TestGetPrefix(t *testing.T) {
    type args struct {
        input string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        // space
        {
            name: "space",
            args: args{
                input: "",
            },
            want: "/",
        },
        // space contained and slash lost
        {
            name: "space contained and slash lost",
            args: args{
                input: "  test1 ",
            },
            want: "test1/",
        },
        // normal
        {
            name: "normal",
            args: args{
                input: "test2",
            },
            want: "test2/",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := GetPrefix(tt.args.input); got != tt.want {
                t.Errorf("GetPrefix() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestJoin(t *testing.T) {
    type args struct {
        sep   string
        input []string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        // space input
        {
            name: "space input",
            args: args{
                sep:   ",",
                input: nil,
            },
            want: "",
        },
        // normal
        {
            name: "space input",
            args: args{
                sep:   ",",
                input: []string{"take", "it"},
            },
            want: "take,it",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Join(tt.args.sep, tt.args.input...); got != tt.want {
                t.Errorf("Join() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestContainforSlice(t *testing.T) {
    type args struct {
        inputSlice   []string
        targetString []string
    }
    tests := []struct {
        name string
        args args
        want bool
    }{
        {
            name: "contain",
            args: args{
                inputSlice:   []string{"alpha", "beta", "gama", "theta"},
                targetString: []string{"beta", "gama"},
            },
            want: true,
        },
        {
            name: "not contain",
            args: args{
                inputSlice:   []string{"alpha", "beta", "gama", "theta"},
                targetString: []string{"beta", "gama", "sigma"},
            },
            want: false,
        },
        {
            name: "wrong order",
            args: args{
                inputSlice:   []string{"alpha", "beta", "gama", "theta"},
                targetString: []string{"gama", "beta"},
            },
            want: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := ContainforSliceInOrder(tt.args.inputSlice, tt.args.targetString...); got != tt.want {
                t.Errorf("ContainforSlice() = %v, want %v", got, tt.want)
            }
        })
    }
}
