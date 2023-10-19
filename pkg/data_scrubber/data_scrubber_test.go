package data_scrubber

import (
	"testing"
)

type ScrubDataTest struct {
	msg             string
	input, expected []byte
}

var scrubDataTests = []ScrubDataTest{
	{"Visa without dashes", []byte("4397945340344828"), []byte("XXXXXXXXXXXX4828")},
	{"Visa with dashes", []byte("4397-9453-4034-4828"), []byte("XXXX-XXXX-XXXX-4828")},
	{"Mastercard without dashes", []byte("5105105105105100"), []byte("XXXXXXXXXXXX5100")},
	{"Mastercard with dashes", []byte("5105-1051-0510-5100"), []byte("XXXX-XXXX-XXXX-5100")},
	{"Discover without dashes", []byte("6011000990139424"), []byte("XXXXXXXXXXXX9424")},
	{"Discover with dashes", []byte("6011-0009-9013-9424"), []byte("XXXX-XXXX-XXXX-9424")},
	{"JCB without dashes", []byte("371449635398431"), []byte("XXXXXXXXXXX8431")},
	{"JCB with dashes", []byte("3714-496353-98431"), []byte("XXXX-XXXXXX-98431")},
	{"Diners Club 35 without dashes", []byte("3566002020360505"), []byte("XXXXXXXXXXXX0505")},
	{"Diners Club 35 with dashes", []byte("3566-0020-2036-0505"), []byte("XXXX-XXXX-XXXX-0505")},
	{"Diners Club 30 without dashes", []byte("30569309025904"), []byte("XXXXXXXXXX5904")},
	{"Diners Club 30 with dashes", []byte("3056-930902-5904"), []byte("XXXX-XXXXXX-5904")},
	{"Social Security Number with dashes", []byte("123-45-6789"), []byte("XXX-XX-X789")},
	{"Social Security Number with spaces", []byte("123 45 6789"), []byte("XXX XX X789")},
	{"Social Security Number without spaces", []byte("123456789"), []byte("XXXXXXX89")},
}

func TestScrubData(t *testing.T) {
	var ds = DataScrubber{}
	ds.Init()
	for _, test := range scrubDataTests {
		input := make([]byte, len(test.input), len(test.input))
		copy(input, test.input)
		result := string(ds.ScrubData(test.input))
		if string(test.expected) != result {
			t.Errorf("test failed: %s\n\tinput: %q\n\texpected: %q\n\tactual: %q\n\n",
				test.msg, input, test.expected, result)
		}
	}
}
