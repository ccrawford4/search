package main

import (
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	documents := [][]byte{
		[]byte(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>CS272 | Welcome</title>
			</head>
			<body>
				<p>Hello World!</p>
				<p>Welcome to <a href="https://cs272-f24.github.io/">CS272</a>!</p>
				<li>
					<a href="https://www.google.com/">Google.com</a>
				</li>
				<div>
					This is a divider
					<a href="/syllabus/">Syllabus</a>
					<a href="/">Home</a>
					<a href="/help/">Help</a>
					<a href="/about/">About</a>
					<a href="https://gobyexample.com">Gobyexample.com</a>
				</div>
			</body>
		</html>
	`),
		[]byte(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Google Home Page</title>
			</head>
			<body>
				<a href="https://www.yahoo.com/" id="321" testId="home-link">Yahoo, Link for the world;!</a>
				<a href="https://www.amazon.com/">!!Amazon!! sink, frank; ; wow</a>
				<a href="mailto:user.lastname@gmail.com">Email-Me!</a>
			</body>
		</html>
	`),
		[]byte(`
			<!DOCTYPE html>
				<head>
					<title>Malformed home page</title>
				</head>
				</body>
				This page is malformed
				<p><a href="https://calum-crawford.com/"></p></a></body>
			</html>
		`),
		[]byte(`
			<!DOCTYPE html>
				<head>
					<title>Testing anchor tag with no url</title>
				<body>
					<a testId="invalid-tag" id="123">https://www.google.com</a>
					<a testId="valid-tag" id="321" href="https://wordle.com/">Wordle</a>
				</body>
				</head>
			</html>
		`),
		[]byte(`
			<!DOCTYPE html>
				 <body>
					<div>Hello</div>
					World
				</body>
			</html>
         `),
		[]byte(`
			<html>
				<head>  
					<title>Style</title>  
					<style>    
						a.blue {     color: blue;    }    
						a.red {      color: red;    }  
					</style>
					<body>  
						<p>Here is a blue link to <a class=\"blue\" href=\"/tests/project01/href.html\">href.html</a></p>  
						<p>And a red link to <a class=\"red\" href=\"/tests/project01/simple.html\">simple.html</a>  
						</p>
					</body>
			</html>
		`),
		[]byte(`
			<html>
				<head>  
					<title>Script</title>  
					<script>    
						for (let i = 0; i < 10; i += 1) {
							console.log("Loop iteration: ", i);
						}
					</script>
					<body>
						<p>Some text</p>
						</p>
					</body>
			</html>
		`),
		nil,
		[]byte(`<!DOCTYPE html>
			<html>
				<style></style>
				<style>body {
					margin-right: 15%;
				}
				</style>
			</html>`),
	}
	expectedWords := [][]string{
		{
			"CS272",
			"|",
			"Welcome",
			"Hello",
			"World",
			"Welcome",
			"to",
			"CS272",
			"Google",
			"com",
			"This",
			"is",
			"a",
			"divider",
			"Syllabus",
			"Home",
			"Help",
			"About",
			"Gobyexample",
			"com",
		},
		{
			"Google",
			"Home",
			"Page",
			"Yahoo",
			"Link",
			"for",
			"the",
			"world",
			"Amazon",
			"sink",
			"frank",
			"wow",
			"Email",
			"Me",
		},
		{
			"Malformed",
			"home",
			"page",
			"This",
			"page",
			"is",
			"malformed",
		},
		{
			"Testing",
			"anchor",
			"tag",
			"with",
			"no",
			"url",
			"https",
			"www",
			"google",
			"com",
			"Wordle",
		},
		{
			"Hello",
			"World",
		},
		{
			"Style",
			"Here",
			"is",
			"a",
			"blue",
			"link",
			"to",
			"href",
			"html",
			"And",
			"a",
			"red",
			"link",
			"to",
			"simple",
			"html",
		},
		{
			"Script",
			"Some",
			"text",
		},
		nil,
		nil,
	}

	expectedHrefs := [][]string{
		{
			"https://cs272-f24.github.io/",
			"https://www.google.com/",
			"/syllabus/",
			"/",
			"/help/",
			"/about/",
			"https://gobyexample.com",
		},
		{
			"https://www.yahoo.com/",
			"https://www.amazon.com/",
			"mailto:user.lastname@gmail.com",
		},
		{
			"https://calum-crawford.com/",
		},
		{
			"https://wordle.com/",
		},
		nil,
		{
			"\\\"/tests/project01/href.html\\\"",
			"\\\"/tests/project01/simple.html\\\"",
		},
		nil,
		nil,
		nil,
	}
	names := []string{
		"simple",
		"punctuation",
		"malformed",
		"no href",
		"no body",
		"style",
		"script",
		"nil",
		"sibling style nodes",
	}
	var tests []struct {
		text                         []byte
		name                         string
		expectedWords, expectedHrefs []string
	}

	for idx := range documents {
		tests = append(tests, struct {
			text                         []byte
			name                         string
			expectedWords, expectedHrefs []string
		}{
			documents[idx],
			names[idx],
			expectedWords[idx],
			expectedHrefs[idx],
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			words, hrefs, _, _ := extract(test.text)
			if !reflect.DeepEqual(words, test.expectedWords) {
				t.Errorf("extract words failed got %q, want %q", words, test.expectedWords)
			}
			if !reflect.DeepEqual(hrefs, test.expectedHrefs) {
				t.Errorf("extract hrefs failed, got %q, want %q", hrefs, test.expectedHrefs)
			}
		})
	}
}
