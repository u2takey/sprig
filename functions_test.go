package sprig

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"text/template"
)

// This is woefully incomplete. Please help.

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestTrimall(t *testing.T) {
	tpl := `{{"$foo$" | trimall "$"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestDefault(t *testing.T) {
	tpl := `{{"" | default "foo"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 234}}`
	if err := runt(tpl, "234"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 2.34}}`
	if err := runt(tpl, "2.34"); err != nil {
		t.Error(err)
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{$v := "foo$bar$baz" | split "$"}}{{$v._0}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

type fixtureTO struct {
	Name, Value string
}

func TestTypeOf(t *testing.T) {
	f := &fixtureTO{"hello", "world"}
	tpl := `{{typeOf .}}`
	if err := runtv(tpl, "*sprig.fixtureTO", f); err != nil {
		t.Error(err)
	}
}

func TestKindOf(t *testing.T) {
	tpl := `{{kindOf .}}`

	f := fixtureTO{"hello", "world"}
	if err := runtv(tpl, "struct", f); err != nil {
		t.Error(err)
	}

	f2 := []string{"hello"}
	if err := runtv(tpl, "slice", f2); err != nil {
		t.Error(err)
	}

	var f3 *fixtureTO = nil
	if err := runtv(tpl, "ptr", f3); err != nil {
		t.Error(err)
	}
}

func TestTypeIs(t *testing.T) {
	f := &fixtureTO{"hello", "world"}
	tpl := `{{if typeIs "*sprig.fixtureTO" .}}t{{else}}f{{end}}`
	if err := runtv(tpl, "t", f); err != nil {
		t.Error(err)
	}

	f2 := "hello"
	if err := runtv(tpl, "f", f2); err != nil {
		t.Error(err)
	}
}
func TestTypeIsLike(t *testing.T) {
	f := "foo"
	tpl := `{{if typeIsLike "string" .}}t{{else}}f{{end}}`
	if err := runtv(tpl, "t", f); err != nil {
		t.Error(err)
	}

	// Now make a pointer. Should still match.
	f2 := &f
	if err := runtv(tpl, "t", f2); err != nil {
		t.Error(err)
	}
}
func TestKindIs(t *testing.T) {
	f := &fixtureTO{"hello", "world"}
	tpl := `{{if kindIs "ptr" .}}t{{else}}f{{end}}`
	if err := runtv(tpl, "t", f); err != nil {
		t.Error(err)
	}
	f2 := "hello"
	if err := runtv(tpl, "f", f2); err != nil {
		t.Error(err)
	}
}

func TestEnv(t *testing.T) {
	os.Setenv("FOO", "bar")
	tpl := `{{env "FOO"}}`
	if err := runt(tpl, "bar"); err != nil {
		t.Error(err)
	}
}

func TestExpandEnv(t *testing.T) {
	os.Setenv("FOO", "bar")
	tpl := `{{expandenv "Hello $FOO"}}`
	if err := runt(tpl, "Hello bar"); err != nil {
		t.Error(err)
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	magicWord := "coffee"
	expect := base64.StdEncoding.EncodeToString([]byte(magicWord))

	if expect == magicWord {
		t.Fatal("Encoder doesn't work.")
	}

	tpl := `{{b64enc "coffee"}}`
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{b64dec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, "")
}
func runtv(tpl, expect string, vars interface{}) error {
	fmap := TxtFuncMap()
	t := template.Must(template.New("test").Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
