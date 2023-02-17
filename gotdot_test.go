package gotdot

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, got, expected interface{}) {
	if got != expected {
		t.Errorf("%v != %v", got, expected)
	}
}

func TestParseEnvFileLine(t *testing.T) {

	t.Run("correct line parsing", func(t *testing.T) {
		expectedKey, expectedValue := "TEST_KEY", "TEST_VALUE"
		gotKey, gotValue := parseEnvFileLine("TEST_KEY=TEST_VALUE")
		assertEqual(t, gotKey, expectedKey)
		assertEqual(t, gotValue, expectedValue)
	})

	t.Run("incorrect line parsing", func(t *testing.T) {
		expectedKey, expectedValue := "", ""
		gotKey, gotValue := parseEnvFileLine("TEST_KEYTEST_VALUE")
		assertEqual(t, gotKey, expectedKey)
		assertEqual(t, gotValue, expectedValue)
	})

	t.Run("double '=' sign line parsing", func(t *testing.T) {
		expectedKey, expectedValue := "TEST_KEY", "TEST=VALUE"
		gotKey, gotValue := parseEnvFileLine("TEST_KEY=TEST=VALUE")
		assertEqual(t, gotKey, expectedKey)
		assertEqual(t, gotValue, expectedValue)
	})

	t.Run("comment trim", func(t *testing.T) {
		expectedKey, expectedValue := "TEST_KEY", "TEST_VALUE"
		gotKey, gotValue := parseEnvFileLine("TEST_KEY=TEST_VALUE # comment")
		assertEqual(t, gotKey, expectedKey)
		assertEqual(t, gotValue, expectedValue)
	})

	t.Run("space trim", func(t *testing.T) {
		expectedKey, expectedValue := "TEST_KEY", "TEST_VALUE"
		gotKey, gotValue := parseEnvFileLine("TEST_KEY = TEST_VALUE")
		assertEqual(t, gotKey, expectedKey)
		assertEqual(t, gotValue, expectedValue)
	})

	t.Run("key to upper case", func(t *testing.T) {
		expectedKey, expectedValue := "TEST_KEY", "TEST_VALUE"
		gotKey, gotValue := parseEnvFileLine("test_key=TEST_VALUE")
		assertEqual(t, gotKey, expectedKey)
		assertEqual(t, gotValue, expectedValue)
	})
}

func TestConfigVar(t *testing.T) {

	var testVar ConfigVar

	t.Run("new ConfigVar", func(t *testing.T) {
		expectedKey, expectedInitValue := "TEST_KEY", "TEST_VALUE"
		testVar = NewConfigVar("TEST_KEY", "TEST_VALUE")
		assertEqual(t, testVar.Key, expectedKey)
		assertEqual(t, testVar.initValue, expectedInitValue)
		assertEqual(t, testVar.castedValue, nil)
	})

	t.Run("cast ConfigVar to int64", func(t *testing.T) {
		expectedCastedValue := int64(32)
		testVar = NewConfigVar("TEST_INT", "32")
		got := testVar.Cast("int").castedValue
		assertEqual(t, got, expectedCastedValue)
	})

	t.Run("cast ConfigVar to false", func(t *testing.T) {
		expectedCastedValue := false
		falseVars := []string{"false", "0", "f", "FALSE", "False"}
		for _, v := range falseVars {
			testVar = NewConfigVar("TEST_FALSE", v)
			got := testVar.Cast("bool").castedValue
			assertEqual(t, got, expectedCastedValue)
		}
	})

	t.Run("cast ConfigVar to true", func(t *testing.T) {
		expectedCastedValue := true
		trueVars := []string{"true", "1", "t", "TRUE", "True"}
		for _, v := range trueVars {
			testVar = NewConfigVar("TEST_FALSE", v)
			got := testVar.Cast("bool").castedValue
			assertEqual(t, got, expectedCastedValue)
		}
	})

	t.Run("import casted ConfigVar to variable", func(t *testing.T) {
		expectedCastedValue := int64(1)
		testVar = NewConfigVar("TEST_IMPORT_WITH_CASTING", "1")
		got := testVar.Cast("int").Import()
		assertEqual(t, got, expectedCastedValue)
	})

	t.Run("import casted ConfigVar to variable", func(t *testing.T) {
		expectedNotCastedValue := "1"
		testVar = NewConfigVar("TEST_IMPORT_NO_CASTING", "1")
		got := testVar.Import()
		assertEqual(t, got, expectedNotCastedValue)
	})
}

func TestGotDot(t *testing.T) {
	t.Run("new GotDot", func(t *testing.T) {
		expectedDotEnvPath, expectedMapping := "./.env", make(map[string]*ConfigVar)
		testConfig := NewGotDot()
		assertEqual(t, testConfig.Path, expectedDotEnvPath)
		assertEqual(t, reflect.TypeOf(testConfig.mapping), reflect.TypeOf(expectedMapping))
	})

	t.Run("change .env path", func(t *testing.T) {
		expectedDotEnvPath := "./.env.local"
		testConfig := NewGotDot()
		testConfig.Path = "./.env.local"
		assertEqual(t, testConfig.Path, expectedDotEnvPath)
	})

	t.Run("create env vars and GoDot with Load", func(t *testing.T) {
		testConfig := NewGotDot()
		testConfig.Path = "./.env.test"
		testConfig.Load()

		assertEqual(t, os.Getenv("TEST_0"), "ZERO")
		assertEqual(t, os.Getenv("TEST_1"), "1")
		assertEqual(t, os.Getenv("TEST_4"), "false")
	})

	t.Run("get ConfigVar with Get", func(t *testing.T) {
		testConfig := NewGotDot()
		testConfig.Path = "./.env.test"
		testConfig.Load()

		expected := "ZERO"
		got, _ := testConfig.Get("TEST_0")
		assertEqual(t, got.Import(), expected)

		expected = "1"
		got, _ = testConfig.Get("TEST_1")
		assertEqual(t, got.Import(), expected)

		expected = "false"
		got, _ = testConfig.Get("TEST_4")
		assertEqual(t, got.Import(), expected)
	})

	t.Run("nonexistent key in ConfigVar with Get", func(t *testing.T) {
		testConfig := NewGotDot()
		testConfig.Path = "./.env.test"
		testConfig.Load()

		absentKey := "absent"
		expected := fmt.Sprintf("Env Variable with key '%s' doesn't exist", absentKey)
		_, goterr := testConfig.Get(absentKey)
		assertEqual(t, goterr.Error(), expected)
	})
}
