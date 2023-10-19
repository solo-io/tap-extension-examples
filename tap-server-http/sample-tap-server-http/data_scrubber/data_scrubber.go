package data_scrubber

import (
	"regexp"
)

type DataScrubber struct {
	regexes   []*regexp.Regexp
	skipChars []byte
}

var (
	// SSN without dashes, e.g. 123456789
	SSN_REGEX_1 = `(?:^|\D)([0-9]{9})(?:\D|$)`
	// SSN with dashes, e.g. 123-45-6789
	SSN_REGEX_2 = `(?:^|\D)([0-9]{3}\-[0-9]{2}\-[0-9]{4})(?:\D|$)`
	// SSN with spaces, e.g. 123 45 6789
	SSN_REGEX_3                 = `(?:^|\D)([0-9]{3}\ [0-9]{2}\ [0-9]{4})(?:\D|$)`
	VISA_REGEX_1                = `(?:^|\D)(4[0-9]{3}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	MASTERCARD_REGEX_1          = `(?:^|\D)(5[1-5][0-9]{2}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	DISCOVER_REGEX_1            = `(?:^|\D)(6011(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	AMEX_REGEX_1                = `(?:^|\D)((?:34|37)[0-9]{2}(?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{5})(?:\D|$)`
	JCB_REGEX_1                 = `(?:^|\D)(3[0-9]{3}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	JCB_REGEX_2                 = `(?:^|\D)((?:2131|1800)[0-9]{11})(?:\D|$)`
	DINERS_CLUB_REGEX_1         = `(?:^|\D)(30[0-5][0-9](?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	DINERS_CLUB_REGEX_2         = `(?:^|\D)((?:36|38)[0-9]{2}(?:\ |\-|)[0-9]{6}(?:\ |\-|)[0-9]{4})(?:\D|$)`
	CREDIT_CARD_TRACKER_REGEX_1 = `([1-9][0-9]{2}\-[0-9]{2}\-[0-9]{4}\^\d)`
	CREDIT_CARD_TRACKER_REGEX_2 = `(?:^|\D)(\%?[Bb]\d{13,19}\^[\-\/\.\w\s]{2,26}\^[0-9][0-9][01][0-9][0-9]{3})`
	CREDIT_CARD_TRACKER_REGEX_3 = `(?:^|\D)(\;\d{13,19}\=(?:\d{3}|)(?:\d{4}|\=))`
)

func (ds *DataScrubber) Init() {
	ds.regexes = []*regexp.Regexp{
		regexp.MustCompile(SSN_REGEX_1),
		regexp.MustCompile(SSN_REGEX_2),
		regexp.MustCompile(SSN_REGEX_3),
		regexp.MustCompile(VISA_REGEX_1),
		regexp.MustCompile(MASTERCARD_REGEX_1),
		regexp.MustCompile(DISCOVER_REGEX_1),
		regexp.MustCompile(AMEX_REGEX_1),
		regexp.MustCompile(JCB_REGEX_1),
		regexp.MustCompile(JCB_REGEX_2),
		regexp.MustCompile(DINERS_CLUB_REGEX_1),
		regexp.MustCompile(DINERS_CLUB_REGEX_2),
		regexp.MustCompile(CREDIT_CARD_TRACKER_REGEX_1),
		regexp.MustCompile(CREDIT_CARD_TRACKER_REGEX_2),
		regexp.MustCompile(CREDIT_CARD_TRACKER_REGEX_3),
	}
	ds.skipChars = []byte("-_ ")
}

// Return a new string with sensitive data removed.
func (ds *DataScrubber) ScrubDataString(stringData string) string {
	return string(ds.ScrubData([]byte(stringData)))
}

// Scrub sensitive data (credit cards, social security numbers, etc.) out of
// bodyData. bodyData is modified in place, so callers should copy it if the
// original data is desired. returns the same bodyData passed in
func (ds *DataScrubber) ScrubData(bodyData []byte) []byte {
	scrubData := func(startIndex, endIndex int) {
		scrubEndIndex := int(float64(endIndex-startIndex)*.7) + startIndex
	replaceChars:
		for i := startIndex; i <= scrubEndIndex; i++ {
			for _, skipChar := range ds.skipChars {
				if bodyData[i] == skipChar {
					continue replaceChars
				}
			}
			bodyData[i] = 'X'
		}
	}
	for _, re := range ds.regexes {
		indices := re.FindAllIndex(bodyData, -1)
		for _, index := range indices {
			scrubData(index[0], index[1])
		}
	}
	return bodyData
}
