package hashstructure

import (
	"strings"
	"testing"
	"time"
)

type DefaultCase struct {
	One, Two interface{}
	Match    bool
}

func (tc *DefaultCase) test(t *testing.T, options ...*HashOptions) {
	var opts *HashOptions
	if len(options) == 1 {
		opts = options[0]
	}
	one, err := Hash(tc.One, opts)
	if err != nil {
		t.Errorf("Failed to hash %#v: %s", tc.One, err)
	}
	two, err := Hash(tc.Two, opts)
	if err != nil {
		t.Errorf("Failed to hash %#v: %s", tc.Two, err)
	}
	// Zero is always wrong
	if one == 0 {
		t.Errorf("zero hash: %#v", tc.One)
	}

	// Compare
	if (one == two) != tc.Match {
		t.Errorf("Equality expected: %#v\nFirst: %#v\nSecond: %#v", tc.Match, tc.One, tc.Two)
	}
}

func TestHash_identity(t *testing.T) {
	cases := []interface{}{
		nil,
		"foo",
		42,
		true,
		false,
		[]string{"foo", "bar"},
		[]interface{}{1, nil, "foo"},
		map[string]string{"foo": "bar"},
		map[interface{}]string{"foo": "bar"},
		map[interface{}]interface{}{"foo": "bar", "bar": 0},
		struct {
			Foo string
			Bar []interface{}
		}{
			Foo: "foo",
			Bar: []interface{}{nil, nil, nil},
		},
		&struct {
			Foo string
			Bar []interface{}
		}{
			Foo: "foo",
			Bar: []interface{}{nil, nil, nil},
		},
	}

	for _, tc := range cases {
		// We run the test 100 times to try to tease out variability
		// in the runtime in terms of ordering.
		valuelist := make([]uint64, 100)
		for i := range valuelist {
			v, err := Hash(tc, nil)
			if err != nil {
				t.Errorf("Error: %s\n\n%#v", err, tc)
			}

			valuelist[i] = v
		}

		// Zero is always wrong
		if valuelist[0] == 0 {
			t.Errorf("zero hash: %#v", tc)
		}

		// Make sure all the values match
		t.Logf("%#v: %d", tc, valuelist[0])
		for i := 1; i < len(valuelist); i++ {
			if valuelist[i] != valuelist[0] {
				t.Errorf("non-matching: %d, %d\n\n%#v", i, 0, tc)
			}
		}
	}
}

func TestHash_equal(t *testing.T) {
	type testFoo struct{ Name string }
	type testBar struct{ Name string }

	now := time.Now()

	cases := []DefaultCase{
		{
			map[string]string{"foo": "bar"},
			map[interface{}]string{"foo": "bar"},
			true,
		},

		{
			map[string]interface{}{"1": "1"},
			map[string]interface{}{"1": "1", "2": "2"},
			false,
		},

		{
			struct{ Fname, Lname string }{"foo", "bar"},
			struct{ Fname, Lname string }{"bar", "foo"},
			false,
		},

		{
			struct{ Lname, Fname string }{"foo", "bar"},
			struct{ Fname, Lname string }{"foo", "bar"},
			false,
		},

		{
			struct{ Lname, Fname string }{"foo", "bar"},
			struct{ Fname, Lname string }{"bar", "foo"},
			false,
		},

		{
			testFoo{"foo"},
			testBar{"foo"},
			false,
		},

		{
			struct {
				Foo        string
				unexported string
			}{
				Foo:        "bar",
				unexported: "baz",
			},
			struct {
				Foo        string
				unexported string
			}{
				Foo:        "bar",
				unexported: "bang",
			},
			true,
		},

		{
			struct {
				Foo        string
				unexported string
				hash       uint64
			}{
				Foo:        "bar",
				unexported: "baz",
				hash:       10,
			},
			struct {
				Foo        string
				unexported string
				hash       uint64
			}{
				Foo:        "bar",
				unexported: "bang",
				hash:       20,
			},
			false,
		},

		{
			struct {
				Foo  string
				hash uint64
			}{
				Foo:  "bar",
				hash: 10,
			},
			struct {
				Foo  string
				hash uint64
			}{
				Foo:  "baz",
				hash: 10,
			},
			true,
		},

		{
			struct {
				testFoo
				Foo string
			}{
				Foo:     "bar",
				testFoo: testFoo{Name: "baz"},
			},
			struct {
				testFoo
				Foo string
			}{
				Foo: "bar",
			},
			true,
		},

		{
			struct {
				Foo string
			}{
				Foo: "bar",
			},
			struct {
				testFoo
				Foo string
			}{
				Foo: "bar",
			},
			true,
		},
		{
			now, // contains monotonic clock
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(),
				now.Minute(), now.Second(), now.Nanosecond(), now.Location()), // does not contain monotonic clock
			true,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

func TestHash_equalIgnore(t *testing.T) {
	type Test1 struct {
		Name string
		UUID string `hash:"ignore"`
	}

	type Test2 struct {
		Name string
		UUID string `hash:"-"`
	}

	type TestTime struct {
		Name string
		Time time.Time `hash:"string"`
	}

	type TestTime2 struct {
		Name string
		Time time.Time
	}

	now := time.Now()
	cases := []DefaultCase{
		{
			Test1{Name: "foo", UUID: "foo"},
			Test1{Name: "foo", UUID: "bar"},
			true,
		},

		{
			Test1{Name: "foo", UUID: "foo"},
			Test1{Name: "foo", UUID: "foo"},
			true,
		},

		{
			Test2{Name: "foo", UUID: "foo"},
			Test2{Name: "foo", UUID: "bar"},
			true,
		},

		{
			Test2{Name: "foo", UUID: "foo"},
			Test2{Name: "foo", UUID: "foo"},
			true,
		},
		{
			TestTime{Name: "foo", Time: now},
			TestTime{Name: "foo", Time: time.Time{}},
			false,
		},
		{
			TestTime{Name: "foo", Time: now},
			TestTime{Name: "foo", Time: now},
			true,
		},
		{
			TestTime2{Name: "foo", Time: now},
			TestTime2{Name: "foo", Time: time.Time{}},
			false,
		},
		{
			TestTime2{Name: "foo", Time: now},
			TestTime2{Name: "foo", Time: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(),
				now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
			},
			true,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

func TestHash_stringTagError(t *testing.T) {
	type Test1 struct {
		Name        string
		BrokenField string `hash:"string"`
	}

	type Test2 struct {
		Name        string
		BustedField int `hash:"string"`
	}

	type Test3 struct {
		Name string
		Time time.Time `hash:"string"`
	}

	cases := []struct {
		Test  interface{}
		Field string
	}{
		{
			Test1{Name: "foo", BrokenField: "bar"},
			"BrokenField",
		},
		{
			Test2{Name: "foo", BustedField: 23},
			"BustedField",
		},
		{
			Test3{Name: "foo", Time: time.Now()},
			"",
		},
	}

	for _, tc := range cases {
		_, err := Hash(tc.Test, nil)
		if err != nil {
			if ens, ok := err.(*ErrNotStringer); ok {
				if ens.Field != tc.Field {
					t.Errorf("did not get expected field %#v: got %s wanted %s", tc.Test, ens.Field, tc.Field)
				}
			} else {
				t.Errorf("unknown error %#v: got %s", tc, err)
			}
		}
	}
}

func TestHash_equalNil(t *testing.T) {
	type Test struct {
		Str   *string
		Int   *int
		Map   map[string]string
		Slice []string
	}

	cases := []struct {
		DefaultCase
		ZeroNil bool
	}{
		{
			DefaultCase: DefaultCase{
				One: Test{
					Str:   nil,
					Int:   nil,
					Map:   nil,
					Slice: nil,
				},
				Two: Test{
					Str:   new(string),
					Int:   new(int),
					Map:   make(map[string]string),
					Slice: make([]string, 0),
				},
				Match: true,
			},
			ZeroNil: true,
		},
		{
			DefaultCase: DefaultCase{One: Test{
				Str:   nil,
				Int:   nil,
				Map:   nil,
				Slice: nil,
			},
				Two: Test{
					Str:   new(string),
					Int:   new(int),
					Map:   make(map[string]string),
					Slice: make([]string, 0),
				},
				Match: false,
			},
			ZeroNil: false,
		},
		{
			DefaultCase: DefaultCase{
				One:   nil,
				Two:   0,
				Match: true,
			},
			ZeroNil: true,
		},
		{
			DefaultCase: DefaultCase{
				One:   nil,
				Two:   0,
				Match: true,
			},
			ZeroNil: false,
		},
	}

	for _, tc := range cases {
		tc.test(t, &HashOptions{ZeroNil: tc.ZeroNil})
	}
}

func TestHash_equalSet(t *testing.T) {
	type Test struct {
		Name    string
		Friends []string `hash:"set"`
	}

	cases := []DefaultCase{
		{
			Test{Name: "foo", Friends: []string{"foo", "bar"}},
			Test{Name: "foo", Friends: []string{"bar", "foo"}},
			true,
		},
		{
			Test{Name: "foo", Friends: []string{"foo", "bar"}},
			Test{Name: "foo", Friends: []string{"foo", "bar"}},
			true,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

func TestHash_equalSlice(t *testing.T) {
	cases := []DefaultCase{
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			true,
		},
		{
			[]int{1, 2, 3},
			[]int{3, 2, 1},
			false,
		},
		{
			[]int{1, 1, 1},
			[]int{1, 1, 1},
			true,
		},
		{
			[]int{1, 1, 1},
			[]int{0, 0, 0},
			false,
		},
		{
			[3]int{1, 2, 3},
			[3]int{1, 2, 3},
			true,
		},
		{
			[3]int{1, 2, 3},
			[3]int{3, 2, 1},
			false,
		},
		{
			[3]int{1, 1, 1},
			[3]int{1, 1, 1},
			true,
		},
		{
			[3]int{1, 1, 1},
			[3]int{0, 0, 0},
			false,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

func TestHash_includable(t *testing.T) {
	cases := []DefaultCase{
		{
			testIncludable{Value: "foo"},
			testIncludable{Value: "foo"},
			true,
		},

		{
			testIncludable{Value: "foo", Ignore: "bar"},
			testIncludable{Value: "foo"},
			true,
		},

		{
			testIncludable{Value: "foo", Ignore: "bar"},
			testIncludable{Value: "bar"},
			false,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

func TestHash_includableMap(t *testing.T) {
	cases := []DefaultCase{
		{
			testIncludableMap{Map: map[string]string{"foo": "bar"}},
			testIncludableMap{Map: map[string]string{"foo": "bar"}},
			true,
		},

		{
			testIncludableMap{Map: map[string]string{"foo": "bar", "ignore": "true"}},
			testIncludableMap{Map: map[string]string{"foo": "bar"}},
			true,
		},

		{
			testIncludableMap{Map: map[string]string{"foo": "bar", "ignore": "true"}},
			testIncludableMap{Map: map[string]string{"bar": "baz"}},
			false,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

func TestHash_hashable(t *testing.T) {
	cases := []DefaultCase{
		{
			testHashable{Value: "foo"},
			&testHashablePointer{Value: "foo"},
			true,
		},

		{
			testHashable{Value: "foo1"},
			&testHashablePointer{Value: "foo2"},
			true,
		},
		{
			testHashable{Value: "foo"},
			&testHashablePointer{Value: "bar"},
			false,
		},
		{
			testHashable{Value: "nofoo"},
			&testHashablePointer{Value: "bar"},
			true,
		},
	}

	for _, tc := range cases {
		tc.test(t)
	}
}

type testIncludable struct {
	Value  string
	Ignore string
}

func (t testIncludable) HashInclude(field string, _ interface{}) (bool, error) {
	return field != "Ignore", nil
}

type testIncludableMap struct {
	Map map[string]string
}

func (t testIncludableMap) HashIncludeMap(field string, k, _ interface{}) (bool, error) {
	if field != "Map" {
		return true, nil
	}

	if s, ok := k.(string); ok && s == "ignore" {
		return false, nil
	}

	return true, nil
}

type testHashable struct {
	Value string
}

func (t testHashable) Hash() uint64 {
	if strings.HasPrefix(t.Value, "foo") {
		return 500
	}
	return 100
}

type testHashablePointer struct {
	Value string
}

func (t *testHashablePointer) Hash() uint64 {
	if strings.HasPrefix(t.Value, "foo") {
		return 500
	}
	return 100
}
