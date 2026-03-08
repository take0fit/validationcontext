package validationcontext_test

import (
	"fmt"
	"testing"

	vc "github.com/take0fit/validationcontext"
	"github.com/take0fit/validationcontext/voauto"
)

type benchSource struct {
	Name  string
	Email string
	Age   int
}

type benchNameVO struct{ value string }
type benchEmailVO struct{ value string }
type benchAgeVO struct{ value int }

type benchTarget struct {
	Name  benchNameVO
	Email benchEmailVO
	Age   benchAgeVO
}

func registerBenchConstructors() {
	voauto.Register("NewName", func(v any, ctx *vc.ValidationContext) any {
		value := v.(string)
		ctx.Required(value, "Name", "", false)
		ctx.ValidateMinLength(value, "Name", 2, "")
		return benchNameVO{value: value}
	})
	voauto.Register("NewEmail", func(v any, ctx *vc.ValidationContext) any {
		value := v.(string)
		ctx.Required(value, "Email", "", false)
		ctx.ValidateEmail(value, "Email", "")
		return benchEmailVO{value: value}
	})
	voauto.Register("NewAge", func(v any, ctx *vc.ValidationContext) any {
		value := v.(int)
		ctx.ValidateMinValue(value, "Age", 18, "")
		ctx.ValidateMaxValue(value, "Age", 120, "")
		return benchAgeVO{value: value}
	})
}

func BenchmarkBindAndValidate_Success(b *testing.B) {
	registerBenchConstructors()
	src := &benchSource{Name: "Taro", Email: "taro@example.com", Age: 30}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := voauto.BindAndValidate[benchTarget](src)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if result.Name.value == "" {
			b.Fatal("unexpected empty result")
		}
	}
}

func BenchmarkBindAndValidate_Failure(b *testing.B) {
	registerBenchConstructors()
	src := &benchSource{Name: "", Email: "invalid", Age: 5}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := voauto.BindAndValidate[benchTarget](src)
		if err == nil {
			b.Fatal("expected validation error")
		}
	}
}

func BenchmarkAddError_Aggregate100(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := vc.NewValidationContext()
		for j := 0; j < 100; j++ {
			ctx.AddError(fmt.Sprintf("Field%d", j), "bench error")
		}
		if err := ctx.AggregateError(); err == nil {
			b.Fatal("expected aggregate error")
		}
	}
}
