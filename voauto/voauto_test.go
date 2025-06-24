package voauto

import (
	"errors"
	"fmt"
	"testing"

	"github.com/take0fit/validationcontext"
)

type dummyVO struct{ v string }

func ctorEcho(val any, _ *validationcontext.ValidationContext) any {
	fmt.Printf("ctorEcho called with: %+v (type: %T)\n", val, val)
	result := dummyVO{v: val.(string)}
	fmt.Printf("ctorEcho returning: %+v\n", result)
	return result
}
func ctorFixed(_ any, _ *validationcontext.ValidationContext) any {
	result := dummyVO{v: "FIXED"}
	fmt.Printf("ctorFixed returning: %+v\n", result)
	return result
}

func reset() {
	mu.Lock()
	defer mu.Unlock()
	constructors = make(map[string]ConstructorFunc)
}

func TestVoauto(t *testing.T) {

	t.Run("RegisterSyncOnce", func(t *testing.T) {
		reset()
		Register("NewDummy", ctorEcho)
		Register("NewDummy", ctorFixed) // the second call overwrites the first one

		type src struct{ Val string }
		type dst struct {
			Val dummyVO `vctag:"NewDummy,Val"`
		}

		result, err := BindAndValidate[dst](&src{Val: "AAA"})
		if err != nil {
			t.Fatalf("BindAndValidate err: %v", err)
		}
		if result.Val.v != "FIXED" { // ctorFixed should have been executed (last registered)
			t.Errorf("want FIXED, got %s", result.Val.v)
		}
	})

	t.Run("BindSuccess", func(t *testing.T) {
		cases := []struct {
			name      string
			register  func()
			src       any
			dest      any
			wantValue string
		}{
			{
				name: "ConventionWithoutTag",
				register: func() {
					reset()
					Register("NewVal", ctorEcho) // field name Val -> NewVal
				},
				src:       &struct{ Val string }{"hello"},
				dest:      &struct{ Val dummyVO }{},
				wantValue: "hello",
			},
			{
				name: "ExplicitTagDifferentField",
				register: func() {
					reset()
					Register("CtorX", ctorEcho)
				},
				src: &struct{ Name string }{"XYZ"},
				dest: &struct {
					Val dummyVO `vctag:"CtorX,Name"`
				}{},
				wantValue: "XYZ",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				c.register()

				// Use BindAndValidate
				if c.name == "ConventionWithoutTag" {
					type destType struct{ Val dummyVO }
					fmt.Printf("Testing with src: %+v\n", c.src)
					result, err := BindAndValidate[destType](c.src)
					if err != nil {
						t.Fatalf("BindAndValidate err: %v", err)
					}
					fmt.Printf("Result: %+v\n", result)
					if result.Val.v != c.wantValue {
						t.Errorf("Val = %s, want %s", result.Val.v, c.wantValue)
					}
				} else {
					type destType struct {
						Val dummyVO `vctag:"CtorX,Name"`
					}
					fmt.Printf("Testing with src: %+v\n", c.src)
					result, err := BindAndValidate[destType](c.src)
					if err != nil {
						t.Fatalf("BindAndValidate err: %v", err)
					}
					fmt.Printf("Result: %+v\n", result)
					if result.Val.v != c.wantValue {
						t.Errorf("Val = %s, want %s", result.Val.v, c.wantValue)
					}
				}
			})
		}
	})

	t.Run("BindErrors", func(t *testing.T) {
		cases := []struct {
			name      string
			src       any
			expectErr string
		}{
			{
				name:      "UnregisteredConstructor",
				src:       &struct{ V string }{"x"},
				expectErr: "constructor not found: NoCtor",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				reset()

				type destType struct {
					V dummyVO `vctag:"NoCtor,V"`
				}

				_, err := BindAndValidate[destType](c.src)
				if err == nil || err.Error() != c.expectErr {
					t.Fatalf("want %q, got %v", c.expectErr, err)
				}
			})
		}
	})

	t.Run("BindAndValidate", func(t *testing.T) {
		reset()
		// constructor that performs a required check
		Register("NewVal", func(v any, vc *validationcontext.ValidationContext) any {
			str := v.(string)
			fmt.Printf("NewVal constructor called with: %q\n", str)
			vc.Required(str, "Val", "required", false)
			result := dummyVO{v: str}
			fmt.Printf("NewVal constructor returning: %+v, has errors: %v\n", result, vc.HasErrors())
			return result
		})

		type req struct{ Val string }
		type dto struct{ Val dummyVO }

		tests := []struct {
			name      string
			in        req
			wantErr   bool
			errIsAgg  bool
			errCount  int
			wantValue string
		}{
			{
				name:      "Valid",
				in:        req{"OK"},
				wantErr:   false,
				wantValue: "OK",
			},
			{
				name:     "InvalidEmpty",
				in:       req{""},
				wantErr:  true,
				errIsAgg: true,
				errCount: 1,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				fmt.Printf("Testing BindAndValidate with: %+v\n", tc.in)
				got, err := BindAndValidate[dto](&tc.in)

				if tc.wantErr {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
					if tc.errIsAgg {
						var agg *validationcontext.ValidationAggregateError
						if !errors.As(err, &agg) {
							t.Fatalf("expect agg error, got %T", err)
						}
						if len(agg.GetMessages()) != tc.errCount {
							t.Errorf("want %d messages, got %d", tc.errCount, len(agg.GetMessages()))
						}
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				fmt.Printf("BindAndValidate result: %+v\n", got)
				if got.Val.v != tc.wantValue {
					t.Errorf("Val = %s, want %s", got.Val.v, tc.wantValue)
				}
			})
		}
	})
}
