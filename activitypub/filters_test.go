package activitypub

import (
	pub "github.com/go-ap/activitypub"
	"reflect"
	"testing"
	"time"
)

func TestFilters_GetLink(t *testing.T) {
	val := pub.IRI("http://example.com")
	f := Filters{
		IRI: val,
	}

	if f.GetLink() != val {
		t.Errorf("Invalid Link returned %s, expected %s", f.GetLink(), val)
	}
}

func TestFilters_IRIs(t *testing.T) {
	val := "http://example.com"
	val1 := "http://example1.com"
	val2 := "http://example1.com/test"
	f := Filters{
		ItemKey: CompStrs{CompStr{Str: val}, CompStr{Str: val1}, CompStr{Str: val2}},
	}
	fullIris := CompStrs{
		CompStr{Str: val},
		CompStr{Str: val1},
		CompStr{Str: val2},
	}

	if !f.IRIs().Contains(CompStr{Str: val}) {
		t.Errorf("Invalid IRIs returned %v, expected %s", f.IRIs(), val)
	}
	if !f.IRIs().Contains(CompStr{Str: val1}) {
		t.Errorf("Invalid IRIs returned %v, expected %s", f.IRIs(), val1)
	}
	if !f.IRIs().Contains(CompStr{Str: val2}) {
		t.Errorf("Invalid IRIs returned %v, expected %s", f.IRIs(), val2)
	}
	if !reflect.DeepEqual(f.IRIs(), fullIris) {
		t.Errorf("Invalid IRIs returned %v, expected %s", f.IRIs(), fullIris)
	}

}

func TestFilters_Page(t *testing.T) {
	t.Skipf("TODO")
}

func TestFilters_Types(t *testing.T) {
	t.Skipf("TODO")
}

func TestFromRequest(t *testing.T) {
	t.Skipf("TODO")
}

func TestHash_String(t *testing.T) {
	t.Skipf("TODO")
}

func TestValidActivityCollection(t *testing.T) {
	t.Skipf("TODO")
}

func mockItem() pub.Object {
	return pub.Object{
		ID:           "",
		Type:         "",
		Name:         nil,
		Attachment:   nil,
		AttributedTo: nil,
		Audience:     nil,
		Content:      nil,
		Context:      nil,
		MediaType:    "",
		EndTime:      time.Time{},
		Generator:    nil,
		Icon:         nil,
		Image:        nil,
		InReplyTo:    nil,
		Location:     nil,
		Preview:      nil,
		Published:    time.Time{},
		Replies:      nil,
		StartTime:    time.Time{},
		Summary:      nil,
		Tag:          nil,
		Updated:      time.Time{},
		URL:          nil,
		To:           nil,
		Bto:          nil,
		CC:           nil,
		BCC:          nil,
		Duration:     0,
		Likes:        nil,
		Shares:       nil,
		Source:       pub.Source{},
	}
}

func EqualsString(s string) CompStr {
	return CompStr{Operator: "=", Str: s}
}
func IRIsFilter(iris ...pub.IRI) CompStrs {
	r := make(CompStrs, len(iris))
	for i, iri := range iris {
		r[i] = EqualsString(iri.String())
	}
	return r
}
func TestFilters_Actors(t *testing.T) {
	f := Filters{
		Actor: &Filters{Key: []Hash{Hash("test")}},
	}

	if f.Actors() == nil {
		t.Errorf("Actors() should not return nil")
		return
	}
	act := mockActivity()
	act.Actor = pub.IRI("/actors/test")
	t.Run("exists", func(t *testing.T) {
		if !testItInIRIs(IRIsFilter(f.Actors()...), act.Actor) {
			t.Errorf("filter %v doesn't contain any of %v", f.Objects(), act.Actor)
		}
	})
	act.Actor = pub.ItemCollection{pub.IRI("/actors/test123"), pub.IRI("https://example.com")}
	t.Run("missing", func(t *testing.T) {
		if testItInIRIs(IRIsFilter(f.Actors()...), act.Actor) {
			t.Errorf("filter %v shouldn't contain any of %v", f.Objects(), act.Actor)
		}
	})
}

func testItInIRIs(iris CompStrs, items ...pub.Item) bool {
	contains := false
	for _, val := range items {
		if val.IsCollection() {
			pub.OnCollectionIntf(val, func(c pub.CollectionInterface) error {
				for _, it := range c.Collection() {
					if filterItem(iris, it) {
						contains = true
						return nil
					}
				}
				return nil
			})
		}
		if filterItemCollections(iris, val) {
			contains = true
			break
		}
	}
	return contains
}

func TestFilters_AttributedTo(t *testing.T) {
	f := Filters{
		InReplTo: CompStrs{CompStr{Str: "test"}},
	}

	if f.InReplyTo() == nil {
		t.Errorf("InReplyTo() should not return nil")
		return
	}
	it := mockItem()
	it.InReplyTo = pub.ItemCollection{pub.IRI("/objects/test")}
	t.Run("exists", func(t *testing.T) {
		if !testItInIRIs(f.InReplyTo(), it.InReplyTo) {
			t.Errorf("filter %v doesn't contain any of %v", f.InReplyTo(), it.InReplyTo)
		}
	})
	it.InReplyTo = pub.ItemCollection{pub.IRI("/objects/test123"), pub.IRI("https://example.com")}
	t.Run("missing", func(t *testing.T) {
		if testItInIRIs(f.InReplyTo(), it.InReplyTo) {
			t.Errorf("filter %v shouldn't contain any of %v", f.InReplyTo(), it.InReplyTo)
		}
	})
}

func TestFilters_Audience(t *testing.T) {
	type args struct {
		filters CompStrs
		valArr  pub.ItemCollection
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "basic-equality",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("ana"),
				},
			},
			want: true,
		},
		{
			name: "basic-equality-with-nil-first",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					nil,
					pub.IRI("ana"),
				},
			},
			want: true,
		},
		{
			name: "basic-like",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("ana"),
				},
			},
			want: true,
		},
		{
			name: "basic-like-with-longer-value",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("anathema"),
				},
			},
			want: true,
		},
		{
			name: "basic-different",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("bob"),
				},
			},
			want: true,
		},
		{
			name: "basic-different-with-empty-values",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					nil,
					pub.IRI(""),
					pub.IRI("bob"),
				},
			},
			want: true,
		},
		{
			name: "basic-false-equality",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("bob"),
				},
			},
			want: false,
		},
		{
			name: "basic-false-equality-with-nil-first",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					nil,
					pub.IRI("bob"),
				},
			},
			want: false,
		},
		{
			name: "basic-false-like",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("bob"),
				},
			},
			want: false,
		},
		{
			name: "basic-false-like-with-longer-value",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("bobsyouruncle"),
				},
			},
			want: false,
		},
		{
			name: "basic-false-different",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					pub.IRI("ana"),
				},
			},
			want: false,
		},
		{
			name: "basic-false-different-with-empty-values",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: pub.ItemCollection{
					nil,
					pub.IRI(""),
				},
			},
			want: true,
		},
		{
			name: "one value: exact match success",
			args: args{
				CompStrs{StringEquals("ana")},
				pub.ItemCollection{
					pub.IRI("ana"),
				},
			},
			want: true,
		},
		{
			name: "one value: exact match failure",
			args: args{
				CompStrs{StringEquals("ana")},
				pub.ItemCollection{
					pub.IRI("na"),
				},
			},
			want: false,
		},
		{
			name: "one value: partial match success",
			args: args{
				CompStrs{StringLike("ana")},
				pub.ItemCollection{
					pub.IRI("analema"),
				},
			},
			want: true,
		},
		{
			name: "one value: exact match failure",
			args: args{
				CompStrs{StringLike("ana")},
				pub.ItemCollection{
					pub.IRI("na"),
				},
			},
			want: false,
		},
		{
			name: "one value: negated match success",
			args: args{
				CompStrs{StringDifferent("ana")},
				pub.ItemCollection{
					pub.IRI("lema"),
				},
			},
			want: true,
		},
		{
			name: "one value: negated match failure",
			args: args{
				CompStrs{StringDifferent("ana")},
				pub.ItemCollection{
					pub.IRI("ana"),
				},
			},
			want: false,
		},
		// multiple filters
		{
			name: "multi filters: exact match success",
			args: args{
				CompStrs{StringEquals("ana")},
				pub.ItemCollection{
					pub.IRI("not-matching"),
					pub.IRI("ana"),
				},
			},
			want: true,
		},
		{
			name: "multi filters: exact match failure",
			args: args{
				CompStrs{StringEquals("ana")},
				pub.ItemCollection{
					pub.IRI("not-matching"),
					pub.IRI("na"),
				},
			},
			want: false,
		},
		{
			name: "multi filters: partial match success",
			args: args{
				CompStrs{StringLike("ana")},
				pub.ItemCollection{
					pub.IRI("not-matching"),
					pub.IRI("analema"),
				},
			},
			want: true,
		},
		{
			name: "multi filters: exact match failure",
			args: args{
				CompStrs{StringLike("ana")},
				pub.ItemCollection{
					pub.IRI("not-matching"),
					pub.IRI("na"),
				},
			},
			want: false,
		},
		{
			name: "multi filters: negated match success",
			args: args{
				CompStrs{StringDifferent("ana")},
				pub.ItemCollection{
					pub.IRI("not-matching"),
					pub.IRI("lema"),
				},
			},
			want: true,
		},
		{
			name: "multi filters: negated match failure",
			args: args{
				CompStrs{StringDifferent("ana")},
				pub.ItemCollection{
					pub.IRI("not-matching"),
					pub.IRI("ana"),
				},
			},
			want: false,
		},
		{
			name: "existing_matching",
			args: args{
				filters: CompStrs{CompStr{Str: "/actors/test"}},
				valArr:  pub.ItemCollection{pub.IRI("/actors/test")},
			},
			want: true,
		},
		{
			name: "existing_not_matching",
			args: args{
				filters: CompStrs{CompStr{Str: "/actors/test"}},
				valArr:  pub.ItemCollection{pub.IRI("/actors/test123"), pub.IRI("https://example.com")},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterAudience(tt.args.filters, tt.args.valArr); got != tt.want {
				t.Errorf("filterAudience() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilters_Context(t *testing.T) {
	f := Filters{
		OP: CompStrs{EqualsString("test")},
	}
	if f.Context() == nil {
		t.Errorf("Context() should not return nil")
		return
	}
	it := mockItem()
	it.Context = pub.IRI("/objects/test")
	t.Run("exists", func(t *testing.T) {
		if !testItInIRIs(f.Context(), it.Context) {
			t.Errorf("filter %v doesn't contain any of %v", f.Context(), it.Context)
		}
	})
	it.Context = pub.ItemCollection{pub.IRI("/objects/test123"), pub.IRI("https://example.com")}
	t.Run("missing", func(t *testing.T) {
		if testItInIRIs(f.Context(), it.Context) {
			t.Errorf("filter %v shouldn't contain any of %v", f.Context(), it.Context)
		}
	})
}

func TestFilters_InReplyTo(t *testing.T) {
	f := Filters{
		InReplTo: CompStrs{EqualsString("test")},
	}
	if f.InReplyTo() == nil {
		t.Errorf("InReplyTo() should not return nil")
		return
	}
	it := mockItem()
	it.InReplyTo = pub.ItemCollection{pub.IRI("/objects/test")}
	t.Run("exists", func(t *testing.T) {
		if !testItInIRIs(f.InReplyTo(), it.InReplyTo) {
			t.Errorf("filter %v doesn't contain any of %v", f.InReplyTo(), it.InReplyTo)
		}
	})
	it.InReplyTo = pub.ItemCollection{pub.IRI("/objects/test123"), pub.IRI("https://example.com")}
	t.Run("missing", func(t *testing.T) {
		if testItInIRIs(f.InReplyTo(), it.InReplyTo) {
			t.Errorf("filter %v shouldn't contain any of %v", f.InReplyTo(), it.InReplyTo)
		}
	})
}

func TestFilters_MediaTypes(t *testing.T) {
	tests := []struct {
		name string
		args Filters
		want []pub.MimeType
	}{
		{
			name: "empty",
			args: Filters{
				MedTypes: []pub.MimeType{},
			},
			want: []pub.MimeType{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.MediaTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filters.MediaTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilters_Names(t *testing.T) {
	tests := []struct {
		name string
		args Filters
		want CompStrs
	}{
		{
			name: "empty",
			args: Filters{
				Name: CompStrs{},
			},
			want: CompStrs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.Names(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filters.Names() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockActivity() pub.Activity {
	return pub.Activity{
		ID:           "",
		Type:         "",
		Name:         nil,
		Attachment:   nil,
		AttributedTo: nil,
		Audience:     nil,
		Content:      nil,
		Context:      nil,
		MediaType:    "",
		EndTime:      time.Time{},
		Generator:    nil,
		Icon:         nil,
		Image:        nil,
		InReplyTo:    nil,
		Location:     nil,
		Preview:      nil,
		Published:    time.Time{},
		Replies:      nil,
		StartTime:    time.Time{},
		Summary:      nil,
		Tag:          nil,
		Updated:      time.Time{},
		URL:          nil,
		To:           nil,
		Bto:          nil,
		CC:           nil,
		BCC:          nil,
		Duration:     0,
		Actor:        nil,
		Target:       nil,
		Result:       nil,
		Origin:       nil,
		Instrument:   nil,
		Object:       nil,
	}

}
func TestFilters_Objects(t *testing.T) {
	f := Filters{
		Object: &Filters{Key: []Hash{Hash("test")}},
	}
	if f.Objects() == nil {
		t.Errorf("Object() should not return nil")
		return
	}
	act := mockActivity()
	act.Object = pub.IRI("/objects/test")
	t.Run("exists", func(t *testing.T) {
		if !testItInIRIs(IRIsFilter(f.Objects()...), act.Object) {
			t.Errorf("filter %v doesn't contain any of %v", f.Objects(), act.Object)
		}
	})
	act.Object = pub.ItemCollection{pub.IRI("/objects/test123"), pub.IRI("https://example.com")}
	t.Run("missing", func(t *testing.T) {
		if testItInIRIs(IRIsFilter(f.Objects()...), act.Object) {
			t.Errorf("filter %v shouldn't contain any of %v", f.Objects(), act.Object)
		}
	})
}

func TestFilters_Targets(t *testing.T) {
	f := Filters{
		Target: &Filters{Key: []Hash{Hash("test")}},
	}
	act := mockActivity()
	act.Target = pub.IRI("/objects/test")
	t.Run("exists", func(t *testing.T) {
		if !testItInIRIs(IRIsFilter(f.Targets()...), act.Target) {
			t.Errorf("filter %v doesn't contain any of %v", f.Targets(), act.Target)
		}
	})
	act.Target = pub.ItemCollection{pub.IRI("/objects/test123"), pub.IRI("https://example.com")}
	t.Run("missing", func(t *testing.T) {
		if testItInIRIs(IRIsFilter(f.Targets()...), act.Target) {
			t.Errorf("filter %v shouldn't contain any of %v", f.Targets(), act.Target)
		}
	})
}
func TestFilters_URLs(t *testing.T) {
	t.Skipf("TODO")
}

func TestFilters_ItemMatches(t *testing.T) {
	t.Skipf("TODO")
}

func TestFilters_FilterCollection(t *testing.T) {
	t.Skipf("TODO")
}

func Test_filterNaturalLanguageValues(t *testing.T) {
	type args struct {
		filters CompStrs
		valArr  []pub.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "basic-equality",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("ana"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "basic-equality-with-nil-first",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					nil,
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("ana"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "basic-like",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("ana"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "basic-like-with-longer-value",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("anathema"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "basic-different",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("bob"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "basic-different-with-empty-values",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					nil,
					{},
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("bob"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "basic-false-equality",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("bob"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "basic-false-equality-with-nil-first",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					nil,
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("bob"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "basic-false-like",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("bob"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "basic-false-like-with-longer-value",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "~",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("bobsyouruncle"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "basic-false-different",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: pub.Content("ana"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "basic-false-different-with-empty-values",
			args: args{
				filters: CompStrs{
					CompStr{
						Operator: "!",
						Str:      "ana",
					},
				},
				valArr: []pub.NaturalLanguageValues{
					nil,
					{},
				},
			},
			want: true,
		},
		{
			name: "one value: exact match success",
			args: args{
				CompStrs{StringEquals("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("ana"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "one value: exact match failure",
			args: args{
				CompStrs{StringEquals("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("na"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "one value: partial match success",
			args: args{
				CompStrs{StringLike("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("analema"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "one value: exact match failure",
			args: args{
				CompStrs{StringLike("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("na"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "one value: negated match success",
			args: args{
				CompStrs{StringDifferent("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("lema"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "one value: negated match failure",
			args: args{
				CompStrs{StringDifferent("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("ana"),
						},
					},
				},
			},
			want: false,
		},
		// multiple filters
		{
			name: "multi filters: exact match success",
			args: args{
				CompStrs{StringEquals("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("not-matching"),
						},
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("ana"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "multi filters: exact match failure",
			args: args{
				CompStrs{StringEquals("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("not-matching"),
						},
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("na"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "multi filters: partial match success",
			args: args{
				CompStrs{StringLike("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("not-matching"),
						},
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("analema"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "multi filters: exact match failure",
			args: args{
				CompStrs{StringLike("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("not-matching"),
						},
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("na"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "multi filters: negated match success",
			args: args{
				CompStrs{StringDifferent("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("not-matching"),
						},
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("lema"),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "multi filters: negated match failure",
			args: args{
				CompStrs{StringDifferent("ana")},
				[]pub.NaturalLanguageValues{
					{
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("not-matching"),
						},
						pub.LangRefValue{
							Ref:   pub.NilLangRef,
							Value: []byte("ana"),
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterNaturalLanguageValues(tt.args.filters, tt.args.valArr...); got != tt.want {
				t.Errorf("filterNaturalLanguageValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
