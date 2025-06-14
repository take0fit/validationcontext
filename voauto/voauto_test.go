package voauto

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/take0fit/validationcontext"
)

/* ──────────────────────────────────────────────
   Test-dummy VO / constructors
   ────────────────────────────────────────────── */

type dummyVO struct{ v string }

func ctorEcho(val any, _ *validationcontext.ValidationContext) any {
	return dummyVO{v: val.(string)}
}
func ctorFixed(_ any, _ *validationcontext.ValidationContext) any {
	return dummyVO{v: "FIXED"}
}

/* ──────────────────────────────────────────────
   Helpers
   ────────────────────────────────────────────── */

func reset() {
	mu.Lock()
	defer mu.Unlock()
	registry = map[string]Constructor{}
	onceMap = map[string]*sync.Once{}
}

/* ──────────────────────────────────────────────
   Tests
   ────────────────────────────────────────────── */

func TestVoauto(t *testing.T) {

	t.Run("RegisterSyncOnce", func(t *testing.T) {
		reset()
		Register("NewDummy", ctorEcho)
		Register("NewDummy", ctorFixed) // the second call should be ignored

		type src struct{ Val string }
		type dst struct {
			Val dummyVO `vctag:"NewDummy,Val"`
		}
		in := src{Val: "AAA"}
		var out dst

		if err := Bind(&out, &in, validationcontext.NewValidationContext()); err != nil {
			t.Fatalf("Bind err: %v", err)
		}
		if out.Val.v != "AAA" { // ctorEcho should have been executed
			t.Errorf("want AAA, got %s", out.Val.v)
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
				src:       struct{ Val string }{"hello"},
				dest:      &struct{ Val dummyVO }{},
				wantValue: "hello",
			},
			{
				name: "ExplicitTagDifferentField",
				register: func() {
					reset()
					Register("CtorX", ctorEcho)
				},
				src: struct{ Name string }{"XYZ"},
				dest: &struct {
					Val dummyVO `vctag:"CtorX,Name"`
				}{},
				wantValue: "XYZ",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				c.register()
				if err := Bind(c.dest, c.src, validationcontext.NewValidationContext()); err != nil {
					t.Fatalf("Bind err: %v", err)
				}
				got := reflect.ValueOf(c.dest).Elem().Field(0).Interface().(dummyVO).v
				if got != c.wantValue {
					t.Errorf("Val = %s, want %s", got, c.wantValue)
				}
			})
		}
	})

	t.Run("BindErrors", func(t *testing.T) {
		cases := []struct {
			name      string
			dest      any
			src       any
			expectErr string
		}{
			{
				name: "UnregisteredConstructor",
				dest: &struct {
					V dummyVO `vctag:"NoCtor,V"`
				}{},
				src:       &struct{ V string }{"x"},
				expectErr: "constructor NoCtor not registered",
			},
			{
				name:      "DestNotPtr",
				dest:      struct{}{},
				src:       struct{}{},
				expectErr: "dto must be pointer to struct",
			},
			{
				name:      "SrcNotStruct",
				dest:      &struct{}{},
				src:       123,
				expectErr: "src must be struct or pointer to struct",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				reset()
				err := Bind(c.dest, c.src, validationcontext.NewValidationContext())
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
			vc.Required(v.(string), "Val", "required", false)
			return dummyVO{v: v.(string)}
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
