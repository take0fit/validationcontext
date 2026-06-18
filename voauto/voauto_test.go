package voauto

import (
	"errors"
	"strings"
	"testing"

	"github.com/take0fit/validationcontext"
)

type dummyVO struct{ v string }

func ctorEcho(val any, _ *validationcontext.ValidationContext) any {
	return dummyVO{v: val.(string)}
}

func ctorFixed(_ any, _ *validationcontext.ValidationContext) any {
	return dummyVO{v: "FIXED"}
}

func TestVoauto(t *testing.T) {
	t.Run("RegisterSyncOnce", func(t *testing.T) {
		Reset()
		Register("NewDummy", ctorEcho)
		Register("NewDummy", ctorFixed)

		type src struct{ Val string }
		type dst struct {
			Val dummyVO `vctag:"NewDummy,Val"`
		}

		result, err := BindAndValidate[dst](&src{Val: "AAA"})
		if err != nil {
			t.Fatalf("BindAndValidate err: %v", err)
		}
		if result.Val.v != "FIXED" {
			t.Errorf("want FIXED, got %s", result.Val.v)
		}
	})

	t.Run("BindSuccess", func(t *testing.T) {
		cases := []struct {
			name      string
			register  func()
			src       any
			wantValue string
		}{
			{
				name: "ConventionWithoutTag",
				register: func() {
					Reset()
					Register("NewVal", ctorEcho)
				},
				src:       &struct{ Val string }{"hello"},
				wantValue: "hello",
			},
			{
				name: "ExplicitTagDifferentField",
				register: func() {
					Reset()
					Register("CtorX", ctorEcho)
				},
				src:       &struct{ Name string }{"XYZ"},
				wantValue: "XYZ",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				c.register()

				if c.name == "ConventionWithoutTag" {
					type destType struct{ Val dummyVO }
					result, err := BindAndValidate[destType](c.src)
					if err != nil {
						t.Fatalf("BindAndValidate err: %v", err)
					}
					if result.Val.v != c.wantValue {
						t.Errorf("Val = %s, want %s", result.Val.v, c.wantValue)
					}
				} else {
					type destType struct {
						Val dummyVO `vctag:"CtorX,Name"`
					}
					result, err := BindAndValidate[destType](c.src)
					if err != nil {
						t.Fatalf("BindAndValidate err: %v", err)
					}
					if result.Val.v != c.wantValue {
						t.Errorf("Val = %s, want %s", result.Val.v, c.wantValue)
					}
				}
			})
		}
	})

	t.Run("BindErrors", func(t *testing.T) {
		t.Run("UnregisteredConstructor", func(t *testing.T) {
			Reset()

			type destType struct {
				V dummyVO `vctag:"NoCtor,V"`
			}

			_, err := BindAndValidate[destType](&struct{ V string }{"x"})
			var bindErr *BindError
			if !errors.As(err, &bindErr) || bindErr.Kind != BindErrorConstructorMissing {
				t.Fatalf("unexpected bind error: %v", err)
			}
			if err == nil || err.Error() != "constructor not found: NoCtor" {
				t.Fatalf("want %q, got %v", "constructor not found: NoCtor", err)
			}
		})

		t.Run("SourceFieldNotFound", func(t *testing.T) {
			Reset()
			Register("CtorX", ctorEcho)

			type destType struct {
				Val dummyVO `vctag:"CtorX,Missing"`
			}

			_, err := BindAndValidate[destType](&struct{ Name string }{"x"})
			var bindErr *BindError
			if !errors.As(err, &bindErr) || bindErr.Kind != BindErrorSourceFieldMissing {
				t.Fatalf("unexpected bind error: %v", err)
			}
			if err == nil || err.Error() != "source field not found: Missing" {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("InvalidSource", func(t *testing.T) {
			Reset()
			type destType struct{ Val dummyVO }

			_, err := BindAndValidate[destType]("not-struct")
			var bindErr *BindError
			if !errors.As(err, &bindErr) || bindErr.Kind != BindErrorSourceInvalid {
				t.Fatalf("unexpected bind error: %v", err)
			}
			if err == nil || !strings.Contains(err.Error(), "source must be a struct or pointer to struct") {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("InvalidTag", func(t *testing.T) {
			Reset()
			type destType struct {
				Val dummyVO `vctag:",Val"`
			}

			_, err := BindAndValidate[destType](&struct{ Val string }{"x"})
			var bindErr *BindError
			if !errors.As(err, &bindErr) || bindErr.Kind != BindErrorInvalidTag {
				t.Fatalf("unexpected bind error: %v", err)
			}
			if err == nil || !strings.Contains(err.Error(), "constructor key is empty") {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("ConstructorPanic", func(t *testing.T) {
			Reset()
			Register("NewVal", func(v any, vc *validationcontext.ValidationContext) any {
				panic("boom")
			})
			type src struct{ Val string }
			type destType struct{ Val dummyVO }

			_, err := BindAndValidate[destType](&src{Val: "x"})
			var bindErr *BindError
			if !errors.As(err, &bindErr) || bindErr.Kind != BindErrorConstructorPanic {
				t.Fatalf("unexpected bind error: %v", err)
			}
			if err == nil || !strings.Contains(err.Error(), "constructor panic") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	})

	t.Run("BindAndValidate", func(t *testing.T) {
		Reset()
		Register("NewVal", func(v any, vc *validationcontext.ValidationContext) any {
			str := v.(string)
			vc.Required(str, "Val", "required", false)
			return dummyVO{v: str}
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
			{name: "Valid", in: req{"OK"}, wantErr: false, wantValue: "OK"},
			{name: "InvalidEmpty", in: req{""}, wantErr: true, errIsAgg: true, errCount: 1},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
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
				if got.Val.v != tc.wantValue {
					t.Errorf("Val = %s, want %s", got.Val.v, tc.wantValue)
				}
			})
		}
	})
}
