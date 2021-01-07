package lesson

import "testing"

const englishHelloPrefix = "Hello, "
func Hello(name string) string {
	if name == "" {
		 name = "World"
	}
	return englishHelloPrefix + name
}

func Repeat(character string) string  {
	var repeated string
	for i:=0; i < 5; i++ {
		repeated += character
	}
	return repeated
}

func Sum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func TestHello(t *testing.T)  {
	numbers := []int {1, 2, 3, 4, 5}
	numbers = append(numbers, 1)
	Sum(numbers)
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got '%q' want '%q'", got, want)
		}
	}
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})
}