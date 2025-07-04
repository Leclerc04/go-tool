package encryption

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
)

var FakeFieldOfStudies = []string{
	"Acounting",
	"Actuarial Science",
	"Aerospace",
	"Agriculture",
	"Anthropology",
	"Architecture",
	"Area Studies",
	"Biological Agricultural Engineering",
	"Biology",
	"Biomedical Engineering",
	"Chemical Engineering",
	"Chemistry",
	"Civil Engineering",
	"Communication",
	"Computer Engineering",
	"Computer Science",
	"Criminology",
	"Economics",
	"Education",
	"Electrical Engineering",
	"English",
	"Entrepreneurship",
	"Environmental Engineering",
	"Environmental Science",
	"Epidemiology",
	"Finance",
	"Financial Engineering",
	"Fine Arts",
	"Geosciences",
	"Health",
	"History",
	"Human Resource Management",
	"Industrial Engineering",
	"Information Systems",
	"International Business",
	"Language",
	"Law",
	"Literature",
	"Management",
	"Marketing",
	"Materials",
	"Mathematics",
	"Mba",
	"Mechanical Engineering",
	"Medicine",
	"Nuclear Engineering",
	"Operations",
	"Pharmacy",
	"Philosophy",
	"Physics",
	"Political Science",
	"Psychology",
	"Public Affairs",
	"Public Health",
	"Public Management Administration",
	"Public Policy Analysis",
	"Religion",
	"Social Work",
	"Sociology",
	"Statistics",
	"Strategy",
	"Supply Chain And Logistics",
	"Urban Planning",
}

var FakeDegrees = []string{
	"Ph.D.",
	"Master",
	"Undergraduate",
}

func stringToInt(s string) int {
	i := 0
	for _, char := range s {
		i += int(char)
	}
	return i
}

func GetFakeProgramName(programID string) string {
	i := stringToInt(programID)
	fos := FakeFieldOfStudies[i%len(FakeFieldOfStudies)]
	degree := FakeDegrees[i%len(FakeDegrees)]
	return fmt.Sprintf("%s in %s", degree, fos)
}

func GetEncryptedProgramInfo(programID, programInfo string) string {
	salt := stringToInt(programID) % 95

	encodeString := ""

	for _, r := range programInfo {
		rint := int(r)
		if rint < 128 {
			encodeString += string(r)
		} else {
			encodeString += "\\u" + strconv.FormatInt(int64(rint), 16)
		}

	}

	reg := regexp.MustCompile(`\\u`)

	encodeString = reg.ReplaceAllString(encodeString, "%u")
	encodeString = base64.StdEncoding.EncodeToString([]byte(encodeString))

	encryptedName := ""

	for _, r := range encodeString {
		rint := ((int(r) + 63 + salt) % 95) + 32
		encryptedName += string(rune(rint))
	}

	return encryptedName
}
