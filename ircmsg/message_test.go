package ircmsg

import (
	"fmt"
	"reflect"
	"testing"
)

type testcode struct {
	raw     string
	message IrcMessage
}
type testcodewithlen struct {
	raw     string
	length  int
	message IrcMessage
}

var decodelentests = []testcodewithlen{
	{":dan-!d@localhost PRIVMSG dan #test :What a cool message\r\n", 20,
		MakeMessage(nil, "dan-!d@localhost", "PR")},
	{"@time=12732;re TEST *\r\n", 512,
		MakeMessage(map[string]string{"time": "12732", "re": ""}, "", "TEST", "*")},
	{"@time=12732;re TEST *\r\n", 512,
		MakeMessage(map[string]string{"time": "12732", "re": ""}, "", "TEST", "*")},
	{":dan- TESTMSG\r\n", 2048,
		MakeMessage(nil, "dan-", "TESTMSG")},
	{":dan- TESTMSG dan \r\n", 12,
		MakeMessage(nil, "dan-", "TESTMS")},
	{"TESTMSG\r\n", 6,
		MakeMessage(nil, "", "TESTMS")},
	{"TESTMSG\r\n", 7,
		MakeMessage(nil, "", "TESTMSG")},
	{"TESTMSG\r\n", 8,
		MakeMessage(nil, "", "TESTMSG")},
	{"TESTMSG\r\n", 9,
		MakeMessage(nil, "", "TESTMSG")},
}

// map[string]string{"time": "12732", "re": ""}
var decodetests = []testcode{
	{":dan-!d@localhost PRIVMSG dan #test :What a cool message\r\n",
		MakeMessage(nil, "dan-!d@localhost", "PRIVMSG", "dan", "#test", "What a cool message")},
	{"@time=2848 :dan-!d@localhost LIST\r\n",
		MakeMessage(map[string]string{"time": "2848"}, "dan-!d@localhost", "LIST")},
	{"@time=2848 LIST\r\n",
		MakeMessage(map[string]string{"time": "2848"}, "", "LIST")},
	{"LIST\r\n",
		MakeMessage(nil, "", "LIST")},
	{"@time=12732;re TEST *a asda:fs :fhye tegh\r\n",
		MakeMessage(map[string]string{"time": "12732", "re": ""}, "", "TEST", "*a", "asda:fs", "fhye tegh")},
	{"@time=12732;re TEST *\r\n",
		MakeMessage(map[string]string{"time": "12732", "re": ""}, "", "TEST", "*")},
	{":dan- TESTMSG\r\n",
		MakeMessage(nil, "dan-", "TESTMSG")},
	{":dan- TESTMSG dan \r\n",
		MakeMessage(nil, "dan-", "TESTMSG", "dan")},
	{"@time=2019-02-28T19:30:01.727Z ping HiThere!\r\n",
		MakeMessage(map[string]string{"time": "2019-02-28T19:30:01.727Z"}, "", "PING", "HiThere!")},
	{"@+draft/test=hi\\nthere PING HiThere!\r\n",
		MakeMessage(map[string]string{"+draft/test": "hi\nthere"}, "", "PING", "HiThere!")},
}

type testparseerror struct {
	raw string
	err error
}

var decodetesterrors = []testparseerror{
	{"\r\n", ErrorLineIsEmpty},
	{"\r\n    ", ErrorLineIsEmpty},
	{"\r\n ", ErrorLineIsEmpty},
	{" \r\n", ErrorLineIsEmpty},
	{" \r\n ", ErrorLineIsEmpty},
	{"     \r\n  ", ErrorLineIsEmpty},
	{"@tags=tesa\r\n", ErrorLineIsEmpty},
	{"@tags=tested  \r\n", ErrorLineIsEmpty},
	{":dan-   \r\n", ErrorLineIsEmpty},
	{":dan-\r\n", ErrorLineIsEmpty},
	{"@tag1=1;tag2=2 :dan \r\n", ErrorLineIsEmpty},
	{"@tag1=1;tag2=2 :dan      \r\n", ErrorLineIsEmpty},
	{"@tag1=1;tag2=2\x00 :dan      \r\n", ErrorLineContainsBadChar},
	{"@tag1=1;tag2=2\x00 :shivaram PRIVMSG #channel  hi\r\n", ErrorLineContainsBadChar},
}

func TestDecode(t *testing.T) {
	for _, pair := range decodelentests {
		ircmsg, err := ParseLine(pair.raw, true, pair.length)
		if err != nil {
			t.Error(
				"For", pair.raw,
				"Failed to parse line:", err,
			)
		}

		if !reflect.DeepEqual(ircmsg, pair.message) {
			t.Error(
				"For", pair.raw,
				"expected", pair.message,
				"got", ircmsg,
			)
		}
	}
	for _, pair := range decodetests {
		ircmsg, err := ParseLine(pair.raw, true, 0)
		if err != nil {
			t.Error(
				"For", pair.raw,
				"Failed to parse line:", err,
			)
		}

		if !reflect.DeepEqual(ircmsg, pair.message) {
			t.Error(
				"For", pair.raw,
				"expected", pair.message,
				"got", ircmsg,
			)
		}
	}
	for _, pair := range decodetesterrors {
		_, err := ParseLine(pair.raw, true, 0)
		if err != pair.err {
			t.Error(
				"For", pair.raw,
				"expected", pair.err,
				"got", err,
			)
		}
	}
}

var encodetests = []testcode{
	{":dan-!d@localhost PRIVMSG dan #test :What a cool message\r\n",
		MakeMessage(nil, "dan-!d@localhost", "PRIVMSG", "dan", "#test", "What a cool message")},
	{"@time=12732 TEST *a asda:fs :fhye tegh\r\n",
		MakeMessage(map[string]string{"time": "12732"}, "", "TEST", "*a", "asda:fs", "fhye tegh")},
	{"@time=12732 TEST *\r\n",
		MakeMessage(map[string]string{"time": "12732"}, "", "TEST", "*")},
	{"@re TEST *\r\n",
		MakeMessage(map[string]string{"re": ""}, "", "TEST", "*")},
}
var encodelentests = []testcodewithlen{
	{":dan-!d@lo\r\n", 12,
		MakeMessage(nil, "dan-!d@localhost", "PRIVMSG", "dan", "#test", "What a cool message")},
	{"@time=12732 TEST *\r\n", 52,
		MakeMessage(map[string]string{"time": "12732"}, "", "TEST", "*")},
	{"@riohwihowihirgowihre TEST *\r\n", 8,
		MakeMessage(map[string]string{"riohwihowihirgowihre": ""}, "", "TEST", "*", "*")},
}

func TestEncode(t *testing.T) {
	for _, pair := range encodetests {
		line, err := pair.message.Line(true, 0)
		if err != nil {
			t.Error(
				"For", pair.raw,
				"Failed to parse line:", err,
			)
		}

		if line != pair.raw {
			t.Error(
				"For", pair.message,
				"expected", pair.raw,
				"got", line,
			)
		}
	}
	for _, pair := range encodelentests {
		line, err := pair.message.Line(true, pair.length)
		if err != nil {
			t.Error(
				"For", pair.raw,
				"Failed to parse line:", err,
			)
		}

		if line != pair.raw {
			t.Error(
				"For", pair.message,
				"expected", pair.raw,
				"got", line,
			)
		}
	}

	// make sure we fail on no command
	msg := MakeMessage(nil, "example.com", "", "*")
	_, err := msg.Line(true, 0)
	if err == nil {
		t.Error(
			"For", "Test Failure 1",
			"expected", "an error",
			"got", err,
		)
	}

	// make sure we fail with params in right way
	msg = MakeMessage(nil, "example.com", "TEST", "*", "t s", "", "Param after empty!")
	_, err = msg.Line(true, 0)
	if err == nil {
		t.Error(
			"For", "Test Failure 2",
			"expected", "an error",
			"got", err,
		)
	}
}

var testMessages = []IrcMessage{
	{
		tags:           map[string]string{"time": "2019-02-27T04:38:57.489Z", "account": "dan-"},
		clientOnlyTags: map[string]string{"+status": "typing"},
		Prefix:         "dan-!~user@example.com",
		Command:        "TAGMSG",
	},
	{
		clientOnlyTags: map[string]string{"+status": "typing"},
		Command:        "PING", // invalid PING command but we don't care
	},
	{
		tags:    map[string]string{"time": "2019-02-27T04:38:57.489Z"},
		Command: "PING", // invalid PING command but we don't care
		Params:  []string{"12345"},
	},
	{
		tags:    map[string]string{"time": "2019-02-27T04:38:57.489Z", "account": "dan-"},
		Prefix:  "dan-!~user@example.com",
		Command: "PRIVMSG",
		Params:  []string{"#ircv3", ":smiley:"},
	},
	{
		tags:    map[string]string{"time": "2019-02-27T04:38:57.489Z", "account": "dan-"},
		Prefix:  "dan-!~user@example.com",
		Command: "PRIVMSG",
		Params:  []string{"#ircv3", "\x01ACTION writes some specs!\x01"},
	},
	{
		Prefix:  "dan-!~user@example.com",
		Command: "PRIVMSG",
		Params:  []string{"#ircv3", ": long trailing command with langue française in it"},
	},
	{
		Prefix:  "dan-!~user@example.com",
		Command: "PRIVMSG",
		Params:  []string{"#ircv3", " : long trailing command with langue française in it "},
	},
	{
		Prefix:  "shivaram",
		Command: "KLINE",
		Params:  []string{"ANDKILL", "24h", "tkadich", "your", "client", "is", "disconnecting", "too", "much"},
	},
	{
		tags:    map[string]string{"time": "2019-02-27T06:01:23.545Z", "draft/msgid": "xjmgr6e4ih7izqu6ehmrtrzscy"},
		Prefix:  "שיברם",
		Command: "PRIVMSG",
		Params:  []string{"ויקם מלך חדש על מצרים אשר לא ידע את יוסף"},
	},
}

func TestEncodeDecode(t *testing.T) {
	for _, message := range testMessages {
		encoded, err := message.Line(false, 0)
		if err != nil {
			t.Errorf("Couldn't encode %v: %v", message, err)
		}
		parsed, err := ParseLine(encoded, true, 0)
		if err != nil {
			t.Errorf("Couldn't re-decode %v: %v", encoded, err)
		}
		if !reflect.DeepEqual(message, parsed) {
			t.Errorf("After encoding and re-parsing, got different messages:\n%v\n%v", message, parsed)
		}
	}
}

func TestErrorLineTooLongGeneration(t *testing.T) {
	message := IrcMessage{
		tags:    map[string]string{"draft/msgid": "SAXV5OYJUr18CNJzdWa1qQ"},
		Prefix:  "shivaram",
		Command: "PRIVMSG",
		Params:  []string{"aaaaaaaaaaaaaaaaaaaaa"},
	}
	_, err := message.LineBytes(true, 0)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 100; i += 1 {
		message.SetTag(fmt.Sprintf("+client-tag-%d", i), "ok")
	}
	line, err := message.LineBytes(true, 0)
	if err != nil {
		t.Error(err)
	}
	if 4096 < len(line) {
		t.Errorf("line is too long: %d", len(line))
	}

	// add excess tag data, pushing us over the limit
	for i := 100; i < 500; i += 1 {
		message.SetTag(fmt.Sprintf("+client-tag-%d", i), "ok")
	}
	line, err = message.LineBytes(true, 0)
	if err != ErrorLineTooLong {
		t.Error(err)
	}

	message.clientOnlyTags = nil
	for i := 0; i < 500; i += 1 {
		message.SetTag(fmt.Sprintf("server-tag-%d", i), "ok")
	}
	line, err = message.LineBytes(true, 0)
	if err != ErrorLineTooLong {
		t.Error(err)
	}

	message.tags = nil
	message.clientOnlyTags = nil
	for i := 0; i < 200; i += 1 {
		message.SetTag(fmt.Sprintf("server-tag-%d", i), "ok")
		message.SetTag(fmt.Sprintf("+client-tag-%d", i), "ok")
	}
	// client cannot send this much tag data:
	line, err = message.LineBytes(true, 0)
	if err != ErrorLineTooLong {
		t.Error(err)
	}
	// but a server can, since the tags are split between client and server budgets:
	line, err = message.LineBytes(false, 0)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	msg := MakeMessage(
		map[string]string{"time": "2019-02-28T08:12:43.480Z", "account": "shivaram"},
		"shivaram_hexchat!~user@irc.darwin.network",
		"PRIVMSG",
		"#darwin", "what's up guys",
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.LineBytes(false, 0)
	}
}

func BenchmarkParse(b *testing.B) {
	line := "@account=shivaram;draft/msgid=dqhkgglocqikjqikbkcdnv5dsq;time=2019-03-01T20:11:21.833Z :shivaram!~shivaram@good-fortune PRIVMSG #darwin :you're an EU citizen, right? it's illegal for you to be here now"
	for i := 0; i < b.N; i++ {
		ParseLine(line, false, 0)
	}
}
